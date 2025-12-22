package main

import (
	"errors"
	"log"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
)

func newUIDevProxy(rawURL string) (http.Handler, error) {
	target, err := parseProxyURL(rawURL)
	if err != nil {
		return nil, err
	}

	proxy := httputil.NewSingleHostReverseProxy(target)
	originalDirector := proxy.Director
	proxy.Director = func(req *http.Request) {
		originalHost := req.Host
		originalScheme := "http"
		if req.TLS != nil {
			originalScheme = "https"
		}
		originalDirector(req)
		req.Header.Set("X-Forwarded-Proto", originalScheme)
		req.Header.Set("X-Forwarded-Host", originalHost)
		if clientIP, _, err := net.SplitHostPort(req.RemoteAddr); err == nil {
			if prior := req.Header.Get("X-Forwarded-For"); prior != "" {
				req.Header.Set("X-Forwarded-For", prior+", "+clientIP)
			} else {
				req.Header.Set("X-Forwarded-For", clientIP)
			}
		}
		req.Host = originalHost
	}
	proxy.ErrorHandler = func(w http.ResponseWriter, r *http.Request, err error) {
		log.Printf("ui proxy: %v", err)
		http.Error(w, "ui dev server unavailable", http.StatusBadGateway)
	}

	return proxy, nil
}

func parseProxyURL(raw string) (*url.URL, error) {
	trimmed := strings.TrimSpace(raw)
	if trimmed == "" {
		return nil, errors.New("empty dev server URL")
	}

	if !strings.Contains(trimmed, "://") {
		trimmed = "http://" + trimmed
	}

	target, err := url.Parse(trimmed)
	if err != nil {
		return nil, err
	}

	if target.Scheme != "http" && target.Scheme != "https" {
		return nil, errors.New("dev server URL must use http or https")
	}

	if target.Host == "" {
		return nil, errors.New("dev server URL missing host")
	}
	return target, nil
}

func proxyTargetURL(raw string) string {
	target, err := parseProxyURL(raw)
	if err != nil {
		return raw
	}
	return target.String()
}
