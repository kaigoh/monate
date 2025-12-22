package moneropay

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"path"
	"strings"
	"time"
)

const defaultTimeout = 15 * time.Second

// Client wraps the MoneroPay merchant endpoints.
type Client struct {
	Endpoint   string
	APIKey     string
	HTTPClient *http.Client
}

func NewClient(endpoint, apiKey string) *Client {
	return &Client{
		Endpoint: strings.TrimRight(endpoint, "/"),
		APIKey:   apiKey,
	}
}

// ReceiveRequest describes the payload for POST /receive.
type ReceiveRequest struct {
	Amount      int64  `json:"amount"`
	Description string `json:"description,omitempty"`
	CallbackURL string `json:"callback_url,omitempty"`
}

// ReceiveResponse contains the data for a new payment request.
type ReceiveResponse struct {
	Address     string    `json:"address"`
	Amount      int64     `json:"amount"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
}

// PaymentState mirrors the callback + GET /receive response payload.
type PaymentState struct {
	Amount       PaymentAmount        `json:"amount"`
	Complete     bool                 `json:"complete"`
	Description  string               `json:"description"`
	CreatedAt    time.Time            `json:"created_at"`
	Transactions []PaymentTransaction `json:"transactions"`
	Transaction  *PaymentTransaction  `json:"transaction"`
}

type PaymentAmount struct {
	Expected int64 `json:"expected"`
	Covered  struct {
		Total    int64 `json:"total"`
		Unlocked int64 `json:"unlocked"`
	} `json:"covered"`
}

type PaymentTransaction struct {
	Amount          int64     `json:"amount"`
	Confirmations   int       `json:"confirmations"`
	DoubleSpendSeen bool      `json:"double_spend_seen"`
	Fee             int64     `json:"fee"`
	Height          int64     `json:"height"`
	Timestamp       time.Time `json:"timestamp"`
	TxHash          string    `json:"tx_hash"`
	UnlockTime      int64     `json:"unlock_time"`
	Locked          bool      `json:"locked"`
}

func (c *Client) httpClient() *http.Client {
	if c.HTTPClient != nil {
		return c.HTTPClient
	}
	return &http.Client{Timeout: defaultTimeout}
}

func (c *Client) request(ctx context.Context, method, relPath string, body any) (*http.Request, error) {
	if c.Endpoint == "" {
		return nil, fmt.Errorf("moneropay endpoint is required")
	}

	var payload io.Reader
	if body != nil {
		encoded, err := json.Marshal(body)
		if err != nil {
			return nil, err
		}
		payload = bytes.NewReader(encoded)
	}

	u, err := url.Parse(c.Endpoint)
	if err != nil {
		return nil, err
	}
	u.Path = path.Join(u.Path, relPath)

	req, err := http.NewRequestWithContext(ctx, method, u.String(), payload)
	if err != nil {
		return nil, err
	}
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	if c.APIKey != "" {
		req.Header.Set("Authorization", "Bearer "+c.APIKey)
		req.Header.Set("X-API-Key", c.APIKey)
	}

	return req, nil
}

// CreateInvoice calls POST /receive.
func (c *Client) CreateInvoice(ctx context.Context, payload ReceiveRequest) (*ReceiveResponse, error) {
	req, err := c.request(ctx, http.MethodPost, "/receive", payload)
	if err != nil {
		return nil, err
	}
	resp, err := c.httpClient().Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 300 {
		body, _ := io.ReadAll(io.LimitReader(resp.Body, 4<<10))
		return nil, fmt.Errorf("moneropay: create invoice failed (%d): %s", resp.StatusCode, body)
	}
	var decoded ReceiveResponse
	if err := json.NewDecoder(resp.Body).Decode(&decoded); err != nil {
		return nil, err
	}
	return &decoded, nil
}

// ManualCheck fetches payment data via GET /receive/:address.
func (c *Client) ManualCheck(ctx context.Context, address string) (*PaymentState, error) {
	if address == "" {
		return nil, fmt.Errorf("address required")
	}
	req, err := c.request(ctx, http.MethodGet, "/receive/"+url.PathEscape(address), nil)
	if err != nil {
		return nil, err
	}
	resp, err := c.httpClient().Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= http.StatusBadRequest {
		body, _ := io.ReadAll(io.LimitReader(resp.Body, 4<<10))
		return nil, fmt.Errorf("moneropay: manual check failed (%d): %s", resp.StatusCode, body)
	}
	var decoded PaymentState
	if err := json.NewDecoder(resp.Body).Decode(&decoded); err != nil {
		return nil, err
	}
	return &decoded, nil
}

// CallbackPayload matches the POST body MoneroPay sends to callback_url.
type CallbackPayload struct {
	PaymentState
}
