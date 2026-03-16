package config

import "fmt"

type DatabaseConfig struct {
	Host            string
	Port            int
	User            string
	Password        string
	DB              string
	MaxIdleConn     int `mapstructure:"max_idle_conn"`
	MaxOpenConn     int `mapstructure:"max_open_conn"`
	ConnMaxLifetime int `mapstructure:"conn_max_lifetime"` // sec
}

// GetDSN returns the Data Source Name (DSN) for database connection
func (dc *DatabaseConfig) GetDSN() string {
	return fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=disable",
		dc.Host,
		dc.User,
		dc.Password,
		dc.DB,
		dc.Port,
	)
}
