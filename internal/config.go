package config

type Config struct {
	DBUrl  string
	APIKey string
}

func LoadConfig() Config {

	// Load from an env

	return Config{
		DBUrl:  "",
		APIKey: "",
	}
}
