package graph

import (
	"github.com/google/uuid"

	"github.com/kaigoh/monate/v2/graph/model"
	"github.com/kaigoh/monate/v2/graph/scalar"
	"github.com/kaigoh/monate/v2/internal/data"
)

func toModelStore(store *data.Store) *model.Store {
	if store == nil {
		return nil
	}

	return &model.Store{
		ID:        store.ID.String(),
		Name:      store.Name,
		Slug:      store.Slug,
		PublicURL: store.PublicURL,
		Theme:     toModelTheme(store.Theme),
		CreatedAt: store.CreatedAt,
		UpdatedAt: store.UpdatedAt,
	}
}

func toModelTheme(theme data.ThemeSettings) *model.ThemeSettings {
	theme = normalizeTheme(theme)
	result := &model.ThemeSettings{
		PrimaryColor:   theme.PrimaryColor,
		AccentColor:    theme.AccentColor,
		ShowFiatAmount: theme.ShowFiatAmount,
		ShowTicker:     theme.ShowTicker,
		PresetAmounts:  make([]int32, 0, len(theme.PresetAmounts)),
	}
	if theme.BackgroundURL != "" {
		result.BackgroundURL = &theme.BackgroundURL
	}
	if theme.LogoURL != "" {
		result.LogoURL = &theme.LogoURL
	}
	if theme.Headline != "" {
		result.Headline = &theme.Headline
	}
	if theme.CustomCopy != "" {
		result.CustomCopy = &theme.CustomCopy
	}
	for _, val := range theme.PresetAmounts {
		result.PresetAmounts = append(result.PresetAmounts, int32(val))
	}
	return result
}

func toModelInvoice(inv *data.Invoice) *model.Invoice {
	if inv == nil {
		return nil
	}

	modelInvoice := &model.Invoice{
		ID:               inv.ID.String(),
		StoreID:          inv.StoreID.String(),
		Description:      stringPtr(inv.Description),
		Reference:        stringPtr(inv.Reference),
		AmountAtomic:     scalar.Long(inv.AmountAtomic),
		ExpectedAmount:   scalar.Long(inv.ExpectedAmount),
		Currency:         inv.Currency,
		FiatAmount:       inv.FiatAmount,
		MoneroPayAddress: inv.MoneroPayAddress,
		Status:           model.InvoiceStatus(inv.Status),
		Complete:         inv.Complete,
		CoveredTotal:     scalar.Long(inv.CoveredTotal),
		CoveredUnlocked:  scalar.Long(inv.CoveredUnlocked),
		ManualCheckCount: int32(inv.ManualCheckCount),
		CreatedAt:        inv.CreatedAt,
		UpdatedAt:        inv.UpdatedAt,
		ResolvedAt:       inv.ResolvedAt,
	}
	if inv.Store.ID != uuid.Nil {
		modelInvoice.Store = toModelStore(&inv.Store)
	}
	return modelInvoice
}

func stringPtr(val string) *string {
	if val == "" {
		return nil
	}
	return &val
}

func applyThemeInput(current data.ThemeSettings, input *model.ThemeSettingsInput) data.ThemeSettings {
	res := current
	if input == nil {
		return normalizeTheme(res)
	}
	if input.PrimaryColor != nil {
		res.PrimaryColor = *input.PrimaryColor
	}
	if input.AccentColor != nil {
		res.AccentColor = *input.AccentColor
	}
	if input.BackgroundURL != nil {
		res.BackgroundURL = *input.BackgroundURL
	}
	if input.LogoURL != nil {
		res.LogoURL = *input.LogoURL
	}
	if input.Headline != nil {
		res.Headline = *input.Headline
	}
	if input.CustomCopy != nil {
		res.CustomCopy = *input.CustomCopy
	}
	if input.ShowFiatAmount != nil {
		res.ShowFiatAmount = *input.ShowFiatAmount
	}
	if input.ShowTicker != nil {
		res.ShowTicker = *input.ShowTicker
	}
	if len(input.PresetAmounts) > 0 {
		res.PresetAmounts = res.PresetAmounts[:0]
		for _, amt := range input.PresetAmounts {
			res.PresetAmounts = append(res.PresetAmounts, int(amt))
		}
	}
	return normalizeTheme(res)
}

func normalizeTheme(theme data.ThemeSettings) data.ThemeSettings {
	if theme.PrimaryColor == "" {
		theme.PrimaryColor = "#ff6600"
	}
	if theme.AccentColor == "" {
		theme.AccentColor = "#181818"
	}
	if theme.PresetAmounts == nil || len(theme.PresetAmounts) == 0 {
		theme.PresetAmounts = []int{5, 10, 25}
	}
	return theme
}
