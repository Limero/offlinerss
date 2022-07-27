package models

type Paths struct {
	Cache string `json:"cache"`
	Urls  string `json:"urls"`
}

type ServerConfig struct {
	Type     string `json:"type"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type ClientConfig struct {
	Type  string `json:"type"`
	Paths Paths  `json:"paths"`
}

type Config struct {
	Server  ServerConfig   `json:"server"`
	Clients []ClientConfig `json:"clients"`
}
