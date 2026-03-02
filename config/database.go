package config

type DatabaseConfig struct {
	URL                 string   `toml:"url" env:"DATABASE_URL"`
	MigrationAttempts   int      `toml:"migration_attempts"`
	MigrationRetryDelay Duration `toml:"migration_retry_delay"`
}
