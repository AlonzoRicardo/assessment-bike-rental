package database

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/BurntSushi/toml"
	_ "github.com/golang-migrate/migrate/v4/source/file" // This import is critical for using the "file" source driver
	_ "github.com/lib/pq"                                // Import the PostgreSQL driver
)

type Config struct {
	Database DatabaseConfig `toml:"database"`
}

type DatabaseConfig struct {
	User     string `toml:"user"`
	Password string `toml:"password"`
	DBName   string `toml:"dbname"`
	Host     string `toml:"host"`
	Port     int    `toml:"port"`
	SSLMode  string `toml:"sslmode"`
}

func LoadConfig(path string) (*Config, error) {
	config := &Config{}

	// Open the config file
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// Parse the config file
	_, err = toml.NewDecoder(file).Decode(config)
	if err != nil {
		return nil, err
	}

	return config, nil
}

// InitDB initializes and returns a database connection using database/sql
func InitDB(config *DatabaseConfig) (*sql.DB, error) {
	// Build the DSN (Data Source Name)
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%d sslmode=%s",
		config.Host, config.User, config.Password, config.DBName, config.Port, config.SSLMode,
	)

	// Connect to the database
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to the database: %w", err)
	}

	// Optionally, you can ping the database to ensure the connection is established
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping the database: %w", err)
	}

	return db, nil
}
