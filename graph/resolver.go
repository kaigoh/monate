package graph

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require
// here.

import (
	"github.com/kaigoh/monate/v2/internal/config"
	"github.com/kaigoh/monate/v2/internal/notifier"
	"gorm.io/gorm"
)

type Resolver struct {
	DB       *gorm.DB
	Config   *config.Config
	Notifier *notifier.InvoiceNotifier
}
