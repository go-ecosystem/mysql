package mysql

import (
	"fmt"
	"time"

	"gorm.io/gorm/logger"
)

// Config mysql configuration
type Config struct {
	// connection string
	Connection *ConnectionConfig
	// the maximum number of open connections to the database
	MaxOpenConns int
	// the maximum number of connections in the idle connection pool
	MaxIdleConns int
	// the maximum amount of time a connection may be reused in second
	ConnMaxLifetime time.Duration
	// log level
	LogLevel logger.LogLevel
}

// ConnectionConfig connection configuration
type ConnectionConfig struct {
	User      string
	Pwd       string
	Host      string
	Port      string
	DBName    string
	Charset   string
	ParseTime bool
	Loc       string
}

// NewConfig new config with default value
func NewConfig(user, pwd, host, port, dbname, charset string, logLevel logger.LogLevel, opts ...ConfigOption) *Config {
	option := decodeConfigOpts(opts)
	return &Config{
		Connection: &ConnectionConfig{
			User:      user,
			Pwd:       pwd,
			Host:      host,
			Port:      port,
			DBName:    dbname,
			Charset:   charset,
			ParseTime: true,
			Loc:       "Local",
		},
		MaxOpenConns:    option.maxOpenConns,
		MaxIdleConns:    option.maxIdleConns,
		ConnMaxLifetime: option.connMaxLifetime,
		LogLevel:        logLevel,
	}
}

// GenDSN generate DSN string
func (c *Config) GenDSN() string {
	result := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=%s&parseTime=%t&loc=%s",
		c.Connection.User,
		c.Connection.Pwd,
		c.Connection.Host,
		c.Connection.Port,
		c.Connection.DBName,
		c.Connection.Charset,
		c.Connection.ParseTime,
		c.Connection.Loc)
	return result
}
