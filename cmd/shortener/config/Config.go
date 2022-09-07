package config

// Config stores the basic configuration options of the Shortener service
type Config struct {
	// ServerAddress - the Port of the app
	ServerAddress string `env:"SERVER_ADDRESS" envDefault:":8080"`
	// BaseURL - the base URL of the app
	BaseURL string `env:"BASE_URL" envDefault:"http://localhost:8080"`
	// FileStoragePath - the path where the file is stored in case a File based storage is used as a Database
	FileStoragePath string `env:"FILE_STORAGE_PATH" envDefault:""`
	// The database connection DSN
	DatabaseAddress string `env:"DATABASE_DSN" envDefault:""`

	EnableHTTPS bool `env:"ENABLE_HTTPS" envDefault:"false"`

	FileConfig string `env:"CONFIG" envDefault:""`
}
