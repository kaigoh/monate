package data

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/kaigoh/monate/v2/internal/moneropay"
)

type InvoiceStatus string

const (
	InvoiceStatusPending  InvoiceStatus = "PENDING"
	InvoiceStatusPaid     InvoiceStatus = "PAID"
	InvoiceStatusExpired  InvoiceStatus = "EXPIRED"
	InvoiceStatusCanceled InvoiceStatus = "CANCELED"
)

// ThemeSettings is stored as JSON but accessed as a native Go struct.
type ThemeSettings struct {
	PrimaryColor   string `json:"primaryColor"`
	AccentColor    string `json:"accentColor"`
	BackgroundURL  string `json:"backgroundUrl"`
	LogoURL        string `json:"logoUrl"`
	Headline       string `json:"headline"`
	CustomCopy     string `json:"customCopy"`
	ShowFiatAmount bool   `json:"showFiatAmount"`
	ShowTicker     bool   `json:"showTicker"`
	PresetAmounts  []int  `json:"presetAmounts"`
}

// Value implements driver.Valuer so ThemeSettings can be stored as JSON.
func (t ThemeSettings) Value() (driver.Value, error) {
	payload, err := json.Marshal(t)
	if err != nil {
		return nil, err
	}
	return string(payload), nil
}

// Scan hydrates ThemeSettings from JSON.
func (t *ThemeSettings) Scan(value interface{}) error {
	if value == nil {
		*t = ThemeSettings{}
		return nil
	}
	var raw []byte
	switch v := value.(type) {
	case []byte:
		raw = v
	case string:
		raw = []byte(v)
	default:
		return errors.New("unsupported theme settings type")
	}
	if len(raw) == 0 {
		*t = ThemeSettings{}
		return nil
	}
	return json.Unmarshal(raw, t)
}

type Store struct {
	ID                uuid.UUID     `gorm:"type:uuid;primaryKey"`
	Name              string        `gorm:"size:255;not null"`
	Slug              string        `gorm:"size:255;uniqueIndex;not null"`
	PublicURL         string        `gorm:"size:2048;not null"`
	MoneroPayEndpoint string        `gorm:"size:2048;not null"`
	MoneroPayAPIKey   string        `gorm:"size:2048;not null"`
	CallbackSecret    string        `gorm:"size:255;uniqueIndex;not null"`
	Theme             ThemeSettings `gorm:"type:json"`
	CreatedAt         time.Time
	UpdatedAt         time.Time
}

func (s *Store) BeforeCreate(tx *gorm.DB) error {
	if s.ID == uuid.Nil {
		s.ID = uuid.New()
	}
	if s.CallbackSecret == "" {
		s.CallbackSecret = uuid.NewString()
	}
	if len(s.Theme.PresetAmounts) == 0 {
		s.Theme.PresetAmounts = []int{5, 10, 50}
	}
	return nil
}

type Invoice struct {
	ID               uuid.UUID     `gorm:"type:uuid;primaryKey"`
	StoreID          uuid.UUID     `gorm:"type:uuid;index"`
	Store            Store         `gorm:"constraint:OnDelete:CASCADE"`
	Description      string        `gorm:"size:1024"`
	Reference        string        `gorm:"size:255"`
	AmountAtomic     int64         `gorm:"not null"`
	ExpectedAmount   int64         `gorm:"not null"`
	Currency         string        `gorm:"size:16;not null"`
	FiatAmount       float64       `gorm:"not null"`
	MoneroPayAddress string        `gorm:"size:256;uniqueIndex"`
	CallbackURL      string        `gorm:"size:2048"`
	Status           InvoiceStatus `gorm:"size:32;not null;default:PENDING"`
	Complete         bool
	CoveredTotal     int64
	CoveredUnlocked  int64
	ManualCheckCount int
	LastManualCheck  *time.Time
	NextCheckAt      time.Time
	WebhookCompleted bool
	CreatedAt        time.Time
	UpdatedAt        time.Time
	ResolvedAt       *time.Time
}

func (i *Invoice) BeforeCreate(tx *gorm.DB) error {
	now := time.Now().UTC()
	if i.ID == uuid.Nil {
		i.ID = uuid.New()
	}
	if i.ExpectedAmount == 0 {
		i.ExpectedAmount = i.AmountAtomic
	}
	if i.Currency == "" {
		i.Currency = "XMR"
	}
	if i.NextCheckAt.IsZero() {
		i.NextCheckAt = now
	}
	return nil
}

// ApplyPaymentState mutates invoice fields based on the latest MoneroPay state.
func (i *Invoice) ApplyPaymentState(state *moneropay.PaymentState) bool {
	if state == nil {
		return false
	}

	changed := false
	if state.Amount.Expected > 0 && i.ExpectedAmount != state.Amount.Expected {
		i.ExpectedAmount = state.Amount.Expected
		changed = true
	}

	if state.Amount.Covered.Total != i.CoveredTotal {
		i.CoveredTotal = state.Amount.Covered.Total
		changed = true
	}
	if state.Amount.Covered.Unlocked != i.CoveredUnlocked {
		i.CoveredUnlocked = state.Amount.Covered.Unlocked
		changed = true
	}

	if state.Description != "" && state.Description != i.Description {
		i.Description = state.Description
		changed = true
	}

	if state.Complete && (!i.Complete || i.Status != InvoiceStatusPaid) {
		i.Complete = true
		now := time.Now().UTC()
		i.Status = InvoiceStatusPaid
		i.ResolvedAt = &now
		changed = true
	} else if !state.Complete && i.Status == InvoiceStatusPaid && !i.Complete {
		i.Status = InvoiceStatusPending
		changed = true
	}

	return changed
}

// TouchManualCheck updates metadata after a manual check run.
func (i *Invoice) TouchManualCheck(now time.Time, scheduleAfter time.Duration) {
	i.ManualCheckCount++
	i.LastManualCheck = &now
	i.NextCheckAt = now.Add(scheduleAfter)
}
