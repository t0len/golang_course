package modules

import "time"

type PostgreConfig struct {
	Host        string
	Port        string
	Username    string
	Password    string
	DBName      string
	SSLMode     string
	ExecTimeout time.Duration
}

type ServerConfig struct {
	Addr   string
	APIKey string
}
