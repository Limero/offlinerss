package models

type ServerConfig struct {
	Type     string `json:"type"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type ClientConfig struct {
	Type string `json:"type"`
}

type Config struct {
	Server  ServerConfig   `json:"server"`
	Clients []ClientConfig `json:"clients"`
}
