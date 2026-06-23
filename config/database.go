package config

type DatabaseConfig struct {
	URL                 URL      `toml:"url" env:"DATABASE_URL"`
	AutoMigrate         bool     `toml:"auto_migrate"`
	MigrationAttempts   int      `toml:"migration_attempts" env:"-"`
	MigrationRetryDelay Duration `toml:"migration_retry_delay" env:"-"`
}
