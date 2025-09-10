package xui

import (
	"encoding/json"
)

// Top-level structure
type AddClientRequestExternalAPI struct {
	Server   int             `json:"server"`
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
	Clients []ClientParam `json:"clients"`
}

type ClientParam struct {
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
	Client
}

type Client struct {
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

type Inbound struct {
	ID             int             `json:"id"`
	Up             int             `json:"up"`
	Down           int             `json:"down"`
	Total          int             `json:"total"`
	Remark         string          `json:"remark"`
	Enable         bool            `json:"enable"`
	ExpiryTime     int64           `json:"expiryTime"`
	ClientStats    interface{}     `json:"clientStats"` // can refine if you know its structure
	Listen         string          `json:"listen"`
	IP             string          `json:"ip"`
	Port           int             `json:"port"`
	Protocol       string          `json:"protocol"`
	Settings       SettingsDecoded `json:"settings"`
	StreamSettings StreamSettings  `json:"streamSettings"`
}

// nested "streamSettings"
type StreamSettings struct {
	Network  string `json:"network"`
	Security string `json:"security"`
}

type CurrentConfigList struct {
	CurrentConfigs []CurrentConfig `json:"current_config_list"`
}

type CurrentConfig struct {
	Inbound      Inbound `json:"inbound"`
	ClientConfig Client  `json:"client_config"`
}
