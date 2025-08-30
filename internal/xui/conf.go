package xui

import (
	"log"
	"fmt"

	"onward-path/config"

	"encoding/json"
	"os"
)

type _Config struct {
    ServerConfigList []serverConfig `json:"server_list"`
}

func newConfig() *_Config {
    return &_Config{
            ServerConfigList: []serverConfig{},
        }
}

func (c *_Config) Load() error {
	if err := c.loadConfig(); err != nil {
		log.Printf("Couldn't load XUI config : %v", err)
		return err
	}
	log.Printf("XUI Config has been loaded successfully!")

	return nil
}

func (c *_Config) loadConfig() error {
	// Read file
	data, err := os.ReadFile(config.XUI_CONFIG)
	if err != nil {
		log.Printf("Couldn't read xui config file: %v", err)
		return err
	}

	// Parse JSON
    var configDTO _ConfigDTO
	if err := json.Unmarshal(data, &configDTO); err != nil {
		log.Printf("Couldn't unmarshal xui config json: %v", err)
		return err
	}
    fmt.Println(fmt.Printf("Main config IP: '%s'", configDTO.ServerConfigList[0].Host))

    for _, _serverConf := range configDTO.ServerConfigList {
        _serverConfig := serverConfig{
            id:           _serverConf.ID,
            host:         _serverConf.Host,
            port:         _serverConf.Port,
            uriPath:      _serverConf.URIPath,
            baseEndpoint: _serverConf.BaseEndpoint,
            adminUser:    _serverConf.AdminUser,
            adminPass:    _serverConf.AdminPass,
        }
        c.ServerConfigList = append(c.ServerConfigList, _serverConfig)
    }

	return nil
}

type _ConfigDTO struct {
    ServerConfigList []serverConfigDTO `json:"server_list"`
}

// internal config (private)
type serverConfig struct {
	id           string
	host         string
	port         int
	uriPath      string
	baseEndpoint string
	adminUser    string
	adminPass    string
}

// exported JSON version (DTO)
type serverConfigDTO struct {
	ID           string `json:"id"`
	Host         string `json:"host"`
	Port         int    `json:"port"`
	URIPath      string `json:"uri_path"`
	BaseEndpoint string `json:"base_endpoint"`
	AdminUser    string `json:"admin_user"`
	AdminPass    string `json:"admin_pass"`
}

// constructor
func newServerConfig() *serverConfig {
	return &serverConfig{}
}

// --- getters ---
func (c *serverConfig) ID() string           { return c.id }
func (c *serverConfig) Host() string         { return c.host }
func (c *serverConfig) Port() int            { return c.port }
func (c *serverConfig) URIPath() string      { return c.uriPath }
func (c *serverConfig) BaseEndpoint() string { return c.baseEndpoint }
func (c *serverConfig) AdminUser() string    { return c.adminUser }
func (c *serverConfig) AdminPass() string    { return c.adminPass }

// --- JSON conversion ---
func (c *serverConfig) ToJSON() ([]byte, error) {
	dto := serverConfigDTO{
		ID:           c.id,
		Host:         c.host,
		Port:         c.port,
		URIPath:      c.uriPath,
		BaseEndpoint: c.baseEndpoint,
		AdminUser:    c.adminUser,
		AdminPass:    c.adminPass,
	}
	return json.MarshalIndent(dto, "", "  ")
}

func FromJSON(data []byte) (*serverConfig, error) {
	var dto serverConfigDTO
	if err := json.Unmarshal(data, &dto); err != nil {
		return nil, err
	}
	return &serverConfig{
		id:           dto.ID,
		host:         dto.Host,
		port:         dto.Port,
		uriPath:      dto.URIPath,
		baseEndpoint: dto.BaseEndpoint,
		adminUser:    dto.AdminUser,
		adminPass:    dto.AdminPass,
	}, nil
}
