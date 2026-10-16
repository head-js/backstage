package config

const (
	Version   = "0.1.0"
	APIPrefix = "/api/v1"
)

type Config struct {
	Host string
	Port int
}

func Default() Config {
	return Config{
		Host: "127.0.0.1",
		Port: 20080,
	}
}
