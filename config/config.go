package config

import "fmt"

type ENV string

const (
	ENV_LOCAL       = "local"
	ENV_TEST        = "test"
	ENV_DEVELOPMENT = "development"
	ENV_PROD        = "production"
)

type Config struct {
	Postgres    Postgres
	Environment string `default:"local" envconfig:"ENVIRONMENT"`
	BuildTag    string
	BuildDate   string
	GitHash     string
}

type Postgres struct {
	DatabaseUsername   string `default:"postgres_root" envconfig:"DB_USER"`
	DatabasePassword   string `default:"1234" envconfig:"DB_PASSWORD"`
	DatabaseName       string `default:"visiondb" envconfig:"DB_NAME"`
	DatabasePort       string `default:"5432" envconfig:"DB_PORT"`
	DatabaseWriterHost string `default:"localhost" envconfig:"DB_WRITER_HOST"`
	DatabaseReaderHost string `default:"localhost" envconfig:"DB_READER_HOST"`
}

func (c Config) ReaderConnStr() string {
	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s", c.Postgres.DatabaseReaderHost, c.Postgres.DatabasePort, c.Postgres.DatabaseUsername, c.Postgres.DatabasePassword, c.Postgres.DatabaseName)
	if c.Environment == ENV_LOCAL || c.Environment == ENV_TEST {
		connStr += " sslmode=disable"
	}
	return connStr
}

func (c Config) WriterConnStr() string {
	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s", c.Postgres.DatabaseWriterHost, c.Postgres.DatabasePort, c.Postgres.DatabaseUsername, c.Postgres.DatabasePassword, c.Postgres.DatabaseName)
	if c.Environment == ENV_LOCAL || c.Environment == ENV_TEST {
		connStr += " sslmode=disable"
	}
	return connStr
}
