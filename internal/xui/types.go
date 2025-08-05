package xui

import (
)

// Top-level structure
type AddClientRequestExternalAPI struct {
	ID       int    `json:"id"`
	Settings SettingsDecoded `json:"settings"` // Will unmarshal again into SettingsDecoded
}

// Top-level structure
type AddClientRequestInternalAPI struct {
	ID       int    `json:"id"`
	Settings string `json:"settings"` // Will unmarshal again into SettingsDecoded
}

// Actual settings structure (embedded in the string)
type SettingsDecoded struct {
	Clients []Client `json:"clients"`
}

type Client struct {
	ID         string `json:"id"`
	Flow       string    `json:"flow"`
	Email      string    `json:"email"`
	LimitIP    int       `json:"limitIp"`
	TotalGB    int       `json:"totalGB"`
	ExpiryTime int       `json:"expiryTime"`
	Enable     bool      `json:"enable"`
	TgID       string    `json:"tgId"`
	SubID      string    `json:"subId"`
	Reset      int       `json:"reset"`
}
