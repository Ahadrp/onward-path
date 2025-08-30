package config

import (
	"log"
)

const (
	CONFIG_PATH  = "/etc/onward-path/onward-path.conf"
	MYSQL_CONFIG = "/etc/onward-path/mysql.conf"
	XUI_CONFIG = "/etc/onward-path/xui.conf"
)

var ()

type Config struct {
}

func New() *Config {
	return &Config{}
}

func (c Config) Load() error {
	log.Println("Config module has been loaded")
	return nil
}

func (c Config) Run() error {
	log.Println("Config module has been run")
	return nil
}
