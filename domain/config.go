package domain

type ServerConfig struct {
	Name     ServerName `json:"name"`
	Username string     `json:"username"`
	Password string     `json:"password"`
	Hostname string     `json:"hostname"`
}

type ClientConfig struct {
	Name ClientName `json:"name"`
}

type Config struct {
	Server  ServerConfig   `json:"server"`
	Clients []ClientConfig `json:"clients"`
}
