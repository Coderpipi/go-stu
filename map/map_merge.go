package mapmerge

import (
	"encoding/json"
	"errors"
)

const (
	ThemeDocTypeTemplate             = "templates"
	ThemeDocTypeConfig               = "config"
	ThemeDocConfigSettingsDataName   = "settings_data.json"
	ThemeDocConfigSettingsSchemaName = "settings_schema.json"
	CheckoutMerchantThemeID          = "checkout-0000-0000-0000-000000000000"
)

type SettingsForm struct{}
type SettingsResult struct {
	Schema   map[string]interface{} `json:"schema"`
	Settings map[string]interface{} `json:"settings"`
}

func parserCheckoutSettings(result *SettingsResult) error {
	settings := result.Settings
	skinSchema, _ := result.Schema["skin_schema"]
	defaultSettings, ok := skinSchema.([]map[string]interface{})
	if !ok {
		return errors.New("skinSchema cannot transfer []map[string]interface{}")
	}
	finalSettings := make(map[string]interface{})
	for k, v := range settings {
		finalSettings[k] = v
	}
	for _, v := range defaultSettings {
		value, ok := v["settings"]
		if !ok {
			continue
		}
		settings, ok := value.([]map[string]interface{})
		if !ok {
			continue
		}
		for _, setting := range settings {
			settingsJson, err := json.Marshal(setting)
			if err != nil {
				return err
			}
			if err := json.Unmarshal(settingsJson, &finalSettings); err != nil {
				return err
			}
		}

	}
	result.Settings = finalSettings
	return nil
}
