package main

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/extension"
	"github.com/99designs/gqlgen/graphql/handler/lru"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/kaigoh/monate/v2/graph"
	"github.com/kaigoh/monate/v2/internal/config"
	"github.com/kaigoh/monate/v2/internal/data"
	"github.com/kaigoh/monate/v2/internal/database"
	"github.com/kaigoh/monate/v2/internal/moneropay"
	"github.com/kaigoh/monate/v2/internal/notifier"
	"github.com/kaigoh/monate/v2/internal/worker"
	"github.com/kaigoh/monate/v2/web"
	"github.com/vektah/gqlparser/v2/ast"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	db, err := database.Open(cfg.DatabaseDriver, cfg.DatabaseDSN)
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}
	if err := db.AutoMigrate(&data.Store{}, &data.Invoice{}); err != nil {
		log.Fatalf("failed to migrate database: %v", err)
	}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	n := notifier.NewInvoiceNotifier()

	resolver := &graph.Resolver{
		DB:       db,
		Config:   cfg,
		Notifier: n,
	}

	srv := handler.New(graph.NewExecutableSchema(graph.Config{Resolvers: resolver}))

	srv.AddTransport(transport.Options{})
	srv.AddTransport(transport.GET{})
	srv.AddTransport(transport.POST{})

	srv.SetQueryCache(lru.New[*ast.QueryDocument](1000))

	srv.Use(extension.Introspection{})
	srv.Use(extension.AutomaticPersistedQuery{
		Cache: lru.New[string](100),
	})

	mux := http.NewServeMux()
	playgroundHandler := playground.Handler("GraphQL playground", "/query")
	mux.Handle("/playground", playgroundHandler)
	mux.Handle("/query", srv)
	mux.Handle("/webhooks/moneropay/", moneroPayWebhookHandler(db, n))

	uiHandler := http.Handler(web.Handler())
	if devURL := strings.TrimSpace(os.Getenv("MONATE_UI_DEV_SERVER_URL")); devURL != "" {
		proxy, err := newUIDevProxy(devURL)
		if err != nil {
			log.Fatalf("failed to init UI dev proxy: %v", err)
		}
		log.Printf("UI: proxying to %s", proxyTargetURL(devURL))
		uiHandler = proxy
	}
	mux.Handle("/", uiHandler)

	httpServer := &http.Server{
		Addr:    ":" + cfg.Port,
		Handler: mux,
	}

	manual := &worker.ManualChecker{
		DB:        db,
		Interval:  cfg.ManualCheckInterval,
		BatchSize: cfg.ManualCheckBatch,
		Notifier:  n,
		Logger:    log.Default(),
	}
	manual.Start(ctx)

	go func() {
		<-ctx.Done()
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		if err := httpServer.Shutdown(shutdownCtx); err != nil {
			log.Printf("http shutdown error: %v", err)
		}
	}()

	log.Printf("connect to http://localhost:%s/playground for GraphQL playground", cfg.Port)
	if err := httpServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Fatalf("server failed: %v", err)
	}
}

func moneroPayWebhookHandler(db *gorm.DB, n *notifier.InvoiceNotifier) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			w.Header().Set("Allow", http.MethodPost)
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		secret := strings.Trim(strings.TrimPrefix(r.URL.Path, "/webhooks/moneropay/"), "/")
		if secret == "" {
			http.Error(w, "missing secret", http.StatusBadRequest)
			return
		}

		invoiceParam := r.URL.Query().Get("invoice")
		if invoiceParam == "" {
			http.Error(w, "missing invoice parameter", http.StatusBadRequest)
			return
		}

		invoiceID, err := uuid.Parse(invoiceParam)
		if err != nil {
			http.Error(w, "invalid invoice id", http.StatusBadRequest)
			return
		}

		defer r.Body.Close()
		var payload moneropay.CallbackPayload
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			http.Error(w, "invalid payload", http.StatusBadRequest)
			return
		}

		var invoice data.Invoice
		if err := db.Preload("Store").First(&invoice, "id = ?", invoiceID).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				http.Error(w, "invoice not found", http.StatusNotFound)
				return
			}
			http.Error(w, "database error", http.StatusInternalServerError)
			return
		}

		if invoice.Store.CallbackSecret != secret {
			http.Error(w, "secret mismatch", http.StatusForbidden)
			return
		}

		changed := invoice.ApplyPaymentState(&payload.PaymentState)
		invoice.WebhookCompleted = payload.Complete

		if err := db.Save(&invoice).Error; err != nil {
			http.Error(w, "database error", http.StatusInternalServerError)
			return
		}

		if changed && n != nil {
			n.Publish(&invoice)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"status":"ok"}`))
	})
}
