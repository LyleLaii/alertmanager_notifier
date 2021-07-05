package db

type DatabaseInfo struct {
	databaseConfig DatabaseConfig `mapstructure:",squash"`
}

type DatabaseConfig struct {
	postgresConfig PostgresConfig `mapstructure:"postgres,omitempty,squash"`
	sqliteConfig SqliteConfig `mapstructure:"sqlite,omitempty,squash"`
}

type SqliteConfig struct {
	DataPath string `mapstructure:"datapath"`
}

type PostgresConfig struct {
	Host     string `mapstructure:"host"`
	Port     string `mapstructure:"port"`
	Dbname   string `mapstructure:"dbname"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	Sslmode  string `mapstructure:"sslmode"`
}