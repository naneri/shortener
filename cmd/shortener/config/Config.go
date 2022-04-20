package config

type Config struct {
	ServerAddress   string `env:"SERVER_ADDRESS" envDefault:":8080"`
	BaseURL         string `env:"BASE_URL" envDefault:"http://localhost:8080"`
	FileStoragePath string `env:"FILE_STORAGE_PATH" envDefault:""`
	DatabaseAddress string `env:"DATABASE_DSN" envDefault:"host=localhost user=postgres password=mysecretpassword dbname=yandex port=5432 sslmode=disable TimeZone=UTC"`
}
