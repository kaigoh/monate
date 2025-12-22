package worker

import (
	"context"
	"log"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/kaigoh/monate/v2/internal/data"
	"github.com/kaigoh/monate/v2/internal/moneropay"
	"github.com/kaigoh/monate/v2/internal/notifier"
)

// ManualChecker polls MoneroPay for invoices that have not yet triggered a callback.
type ManualChecker struct {
	DB        *gorm.DB
	Interval  time.Duration
	BatchSize int
	Notifier  *notifier.InvoiceNotifier
	Logger    *log.Logger
}

func (m *ManualChecker) Start(ctx context.Context) {
	if m.DB == nil || m.Interval <= 0 {
		return
	}
	go m.runOnce(ctx)
	ticker := time.NewTicker(m.Interval)
	go func() {
		defer ticker.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				m.runOnce(ctx)
			}
		}
	}()
}

func (m *ManualChecker) runOnce(ctx context.Context) {
	now := time.Now().UTC()
	var invoices []data.Invoice
	err := m.DB.
		Preload("Store").
		Where("status = ? AND next_check_at <= ?", data.InvoiceStatusPending, now).
		Order("next_check_at ASC").
		Limit(m.BatchSize).
		Find(&invoices).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		m.logf("manual check query failed: %v", err)
		return
	}
	for _, inv := range invoices {
		m.checkInvoice(ctx, &inv)
	}
}

func (m *ManualChecker) checkInvoice(ctx context.Context, inv *data.Invoice) {
	if inv == nil || inv.StoreID == uuid.Nil || inv.MoneroPayAddress == "" {
		return
	}
	if inv.Store.MoneroPayEndpoint == "" || inv.Store.MoneroPayAPIKey == "" {
		return
	}

	client := moneropay.NewClient(inv.Store.MoneroPayEndpoint, inv.Store.MoneroPayAPIKey)
	childCtx, cancel := context.WithTimeout(ctx, 20*time.Second)
	defer cancel()

	state, err := client.ManualCheck(childCtx, inv.MoneroPayAddress)
	if err != nil {
		m.logf("manual check for invoice %s failed: %v", inv.ID, err)
		inv.TouchManualCheck(time.Now().UTC(), m.Interval)
		if saveErr := m.DB.Model(inv).Updates(map[string]any{
			"manual_check_count": inv.ManualCheckCount,
			"last_manual_check":  inv.LastManualCheck,
			"next_check_at":      inv.NextCheckAt,
		}).Error; saveErr != nil {
			m.logf("failed to persist manual check metadata for %s: %v", inv.ID, saveErr)
		}
		return
	}

	now := time.Now().UTC()
	inv.TouchManualCheck(now, m.Interval)
	changed := inv.ApplyPaymentState(state)

	if err := m.DB.Save(inv).Error; err != nil {
		m.logf("failed to persist manual check invoice %s: %v", inv.ID, err)
		return
	}

	if changed && m.Notifier != nil {
		m.Notifier.Publish(inv)
	}
}

func (m *ManualChecker) logf(format string, args ...any) {
	if m == nil {
		return
	}
	if m.Logger != nil {
		m.Logger.Printf(format, args...)
		return
	}
	log.Printf(format, args...)
}
