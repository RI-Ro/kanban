package config

type Config struct {
	SERVER_IPADDRESS string
	SERVER_PORT      string
	DB_NAME          string
	DB_IPADDRESS     string
	DB_PORT          string
	DB_USER          string
	Development      bool
}

func CreateConfig() Config {
	config := Config{"0.0.0.0", "8080", "tasks", "192.168.10.96", "5432", "postgres", false}
	return config
}
