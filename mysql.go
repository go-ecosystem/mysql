package mysql

import (
	"database/sql"

	"github.com/DATA-DOG/go-sqlmock"

	"log"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
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
	// unregister existed one
	UnregisterByKey(key)

	dsn := config.GenDSN()

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		log.Println("Failed to connect to mysql", dsn, err)
		panic(err)
	}

	// https://gorm.io/docs/connecting_to_the_database.html#Connection-Pool
	sqlDB, err := db.DB()
	if err != nil {
		log.Println(err)
	} else {
		sqlDB.SetConnMaxLifetime(config.ConnMaxLifetime * time.Second)
		sqlDB.SetMaxOpenConns(config.MaxOpenConns)
		sqlDB.SetMaxIdleConns(config.MaxIdleConns)
	}

	// retain the db
	container[key] = db
}

// Unregister unregister default examples
func Unregister() {
	UnregisterByKey(defaultKey)
}

// UnregisterByKey unregister examples by key
func UnregisterByKey(key string) {
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

// MockDB mock default DB
func MockDB() (*sql.DB, sqlmock.Sqlmock) {
	return MockDBByKey(defaultKey)
}

// MockDBByKey mock DB by key
func MockDBByKey(key string) (*sql.DB, sqlmock.Sqlmock) {
	sqlDB, mock, err := sqlmock.New()
	if err != nil {
		panic(err)
	}

	db, err := gorm.Open(mysql.New(mysql.Config{
		Conn: sqlDB,
	}), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		panic(err)
	}

	container[key] = db
	return sqlDB, mock
}
