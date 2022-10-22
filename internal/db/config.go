package db

import (
	"fmt"
)

type Config struct {
	Name       string
	Host       string
	Port       uint16
	User       string
	Password   string
	MaxConns   uint
	IdleConns  uint
	Driver     string
	Migrations *MigrationConfig
}

type MigrationConfig struct {
	Dialect string
	Verbose bool
}

func (cfg *Config) DSN() string {
	return fmt.Sprintf("%v:%v@tcp(%v:%v)/%v?parseTime=true&tx_isolation='%s'&multiStatements=true&maxAllowedPacket=%d",
		cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.Name, txIsolation, maxAllowedPacket)
}
