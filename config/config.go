package config

type (
	ServerConfig struct {
		HttpServer
	}

	HttpServer struct {
		Host string
		Port uint
	}
)

func New() *ServerConfig {
	cfg := &ServerConfig{
		HttpServer: HttpServer{
			Host: "",
			Port: 8080,
		},
	}
	return cfg
}
