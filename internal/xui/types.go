package xui

import (
	"encoding/json"
)

// Top-level structure
type AddClientRequestExternalAPI struct {
	ID       int             `json:"id"`
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
	Flow       string `json:"flow"`
	Email      string `json:"email"`
	LimitIP    int    `json:"limitIp"`
	TotalGB    int    `json:"totalGB"`
	ExpiryTime int    `json:"expiryTime"`
	Enable     bool   `json:"enable"`
	TgID       string `json:"tgId"`
	SubID      string `json:"subId"`
	Reset      int    `json:"reset"`
}

type XUIResponse struct {
	Success bool            `json:"success"`
	Message string          `json:"msg"`
	Obj     json.RawMessage `json:"obj"`
}

type GetClientResponse struct {
	ID        int    `json:"id"`
	InboundID int    `json:"inboundId"`
	Enable    bool   `json:"enable"`
	Email     string `json:"email"`
	Up        int    `json:"up"`
	Down      int    `json:"down"`
	Expiry    int    `json:"expiryTime"`
	Total     int    `json:"total"`
	Reset     int    `json:"reset"`
}
