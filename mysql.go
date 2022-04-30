package mysql

import (
	"database/sql"
	"os"

	"github.com/DATA-DOG/go-sqlmock"

	"log"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

var container = make(map[string]*gorm.DB)

const defaultKey = "default"

// GetDB get default db
func GetDB() *gorm.DB {
	return GetDBByKey(defaultKey)
}

// GetDBByKey get by with key
func GetDBByKey(key string) *gorm.DB {
	return container[key]
}

// Close closes current db connection
func Close() {
	for id := range container {
		db := container[id]
		sqlDB, err := db.DB()
		if err != nil {
			log.Print(err)
			return
		}
		if err := sqlDB.Close(); err != nil {
			log.Print(err)
		}
	}
}

// Register register default examples
func Register(config *Config) {
	RegisterByKey(config, defaultKey)
}

// RegisterByKey register examples by key
func RegisterByKey(config *Config, key string) {
	// deregister existed one
	DeregisterByKey(key)

	dsn := config.GenDSN()

	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold: time.Second,
			LogLevel:      config.LogLevel,
			Colorful:      false,
		},
	)

	const defaultStringSize = 256
	db, err := gorm.Open(mysql.New(mysql.Config{
		DSN:               dsn,
		DefaultStringSize: defaultStringSize,
	}), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			TablePrefix:   config.TablePrefix,
			SingularTable: config.SingularTable,
		},
		Logger: newLogger,
	})
	if err != nil {
		log.Println("Failed to connect to mysql", dsn, err)
		panic(err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		log.Println("get database pool fail,  ", err)
		panic(err)
	}

	// https://gorm.io/docs/connecting_to_the_database.html#Connection-Pool
	sqlDB.SetConnMaxLifetime(config.ConnMaxLifetime * time.Second)
	sqlDB.SetMaxOpenConns(config.MaxOpenConns)
	sqlDB.SetMaxIdleConns(config.MaxIdleConns)

	// retain the db
	container[key] = db
}

// Deregister deregister default examples
func Deregister() {
	DeregisterByKey(defaultKey)
}

// DeregisterByKey deregister examples by key
func DeregisterByKey(key string) {
	db := container[key]
	if db == nil {
		return
	}
	sqlDB, err := db.DB()
	if err != nil {
		log.Print(err)
		return
	}
	if err := sqlDB.Close(); err != nil {
		log.Print(err)
	}
	delete(container, key)
}

// SetLogger replace default logger
func SetLogger(db *gorm.DB, logger logger.Interface) {
	db.Logger = logger
}

// MockDB mock default DB
func MockDB(config *Config) (*sql.DB, sqlmock.Sqlmock) {
	return MockDBByKey(config, defaultKey)
}

// MockDBByKey mock DB by key
func MockDBByKey(config *Config, key string) (*sql.DB, sqlmock.Sqlmock) {
	db, mock, err := sqlmock.New()
	if err != nil {
		panic(err)
	}
	if config == nil {
		config = &Config{}
	}

	dialer := mysql.New(mysql.Config{
		Conn:                      db,
		SkipInitializeWithVersion: true,
	})

	gdb, err := gorm.Open(dialer, &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			TablePrefix:   config.TablePrefix,
			SingularTable: config.SingularTable,
		},
	})

	if err != nil {
		log.Println("connect to database fail, ", err)
		panic(err)
	}

	container[key] = gdb
	return db, mock
}
