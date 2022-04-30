package mysql

import (
	"reflect"
	"testing"
	"time"

	"gorm.io/gorm/logger"
)

func TestNewConfig(t *testing.T) {
	user := "1"
	pwd := "2"
	host := "3"
	port := "4"
	dbname := "5"
	charset := "6"
	logLevel := logger.Info
	tablePrefix := "t_"

	t.Run("without option", func(t *testing.T) {
		got := NewConfig(user, pwd, host, port, dbname, charset, logLevel)

		want := &Config{
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
			MaxOpenConns:    defaultConfigOption().maxOpenConns,
			MaxIdleConns:    defaultConfigOption().maxIdleConns,
			ConnMaxLifetime: defaultConfigOption().connMaxLifetime,
			LogLevel:        logLevel,
		}

		if !reflect.DeepEqual(got, want) {
			t.Errorf("NewConfig() = %v, want %v", got, want)
		}
	})

	t.Run("with maxOpenConns", func(t *testing.T) {
		maxOpenConns := 1000

		got := NewConfig(user, pwd, host, port, dbname, charset, logLevel, WithMaxOpenConns(maxOpenConns))

		want := &Config{
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
			MaxOpenConns:    maxOpenConns,
			MaxIdleConns:    defaultConfigOption().maxIdleConns,
			ConnMaxLifetime: defaultConfigOption().connMaxLifetime,
			LogLevel:        logLevel,
		}

		if !reflect.DeepEqual(got, want) {
			t.Errorf("NewConfig() = %v, want %v", got, want)
		}
	})

	t.Run("with maxIdleConns", func(t *testing.T) {
		maxIdleConns := 2000

		got := NewConfig(user, pwd, host, port, dbname, charset, logLevel, WithMaxIdleConns(maxIdleConns))

		want := &Config{
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
			MaxOpenConns:    defaultConfigOption().maxOpenConns,
			MaxIdleConns:    maxIdleConns,
			ConnMaxLifetime: defaultConfigOption().connMaxLifetime,
			LogLevel:        logLevel,
		}

		if !reflect.DeepEqual(got, want) {
			t.Errorf("NewConfig() = %v, want %v", got, want)
		}
	})

	t.Run("with connMaxLifetime", func(t *testing.T) {
		var connMaxLifetime time.Duration = 5000

		got := NewConfig(user, pwd, host, port, dbname, charset, logLevel, WithConnMaxLifetime(connMaxLifetime))

		want := &Config{
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
			MaxOpenConns:    defaultConfigOption().maxOpenConns,
			MaxIdleConns:    defaultConfigOption().maxIdleConns,
			ConnMaxLifetime: connMaxLifetime,
			LogLevel:        logLevel,
		}

		if !reflect.DeepEqual(got, want) {
			t.Errorf("NewConfig() = %v, want %v", got, want)
		}
	})

	t.Run("with singularTable", func(t *testing.T) {
		singularTable := true

		got := NewConfig(user, pwd, host, port, dbname, charset, logLevel, WithSingularTable(singularTable))

		want := &Config{
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
			MaxOpenConns:    defaultConfigOption().maxOpenConns,
			MaxIdleConns:    defaultConfigOption().maxIdleConns,
			ConnMaxLifetime: defaultConfigOption().connMaxLifetime,
			TablePrefix:     defaultConfigOption().tablePrefix,
			SingularTable:   singularTable,
			LogLevel:        logLevel,
		}

		if !reflect.DeepEqual(got, want) {
			t.Errorf("NewConfig() = %v, want %v", got, want)
		}
	})

	t.Run("with tablePrefix", func(t *testing.T) {
		got := NewConfig(user, pwd, host, port, dbname, charset, logLevel, WithTablePrefix(tablePrefix))

		want := &Config{
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
			MaxOpenConns:    defaultConfigOption().maxOpenConns,
			MaxIdleConns:    defaultConfigOption().maxIdleConns,
			ConnMaxLifetime: defaultConfigOption().connMaxLifetime,
			TablePrefix:     tablePrefix,
			SingularTable:   defaultConfigOption().singularTable,
			LogLevel:        logLevel,
		}

		if !reflect.DeepEqual(got, want) {
			t.Errorf("NewConfig() = %v, want %v", got, want)
		}
	})

	t.Run("with options", func(t *testing.T) {
		maxOpenConns := 1000
		maxIdleConns := 2000
		var connMaxLifetime time.Duration = 5000
		singularTable := true

		var options = make([]ConfigOption, 0)
		options = append(options, WithConnMaxLifetime(connMaxLifetime))
		options = append(options, WithMaxIdleConns(maxIdleConns))
		options = append(options, WithMaxOpenConns(maxOpenConns))
		options = append(options, WithSingularTable(singularTable))
		options = append(options, WithTablePrefix(tablePrefix))

		got := NewConfig(user, pwd, host, port, dbname, charset, logLevel, options...)

		want := &Config{
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
			MaxOpenConns:    maxOpenConns,
			MaxIdleConns:    maxIdleConns,
			ConnMaxLifetime: connMaxLifetime,
			LogLevel:        logLevel,
			TablePrefix:     tablePrefix,
			SingularTable:   singularTable,
		}

		if !reflect.DeepEqual(got, want) {
			t.Errorf("NewConfig() = %v, want %v", got, want)
		}
	})
}

func TestConfig_GenDSN(t *testing.T) {
	type fields struct {
		Connection      *ConnectionConfig
		MaxOpenConns    int
		MaxIdleConns    int
		ConnMaxLifetime time.Duration
		LogLevel        logger.LogLevel
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "",
			fields: fields{
				Connection: &ConnectionConfig{
					User:      "1",
					Pwd:       "2",
					Host:      "3",
					Port:      "4",
					DBName:    "5",
					Charset:   "6",
					ParseTime: true,
					Loc:       "Local",
				},
			},
			want: "1:2@tcp(3:4)/5?charset=6&parseTime=true&loc=Local",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Config{
				Connection:      tt.fields.Connection,
				MaxOpenConns:    tt.fields.MaxOpenConns,
				MaxIdleConns:    tt.fields.MaxIdleConns,
				ConnMaxLifetime: tt.fields.ConnMaxLifetime,
				LogLevel:        tt.fields.LogLevel,
			}
			if got := c.GenDSN(); got != tt.want {
				t.Errorf("Config.GenDSN() = %v, want %v", got, tt.want)
			}
		})
	}
}
