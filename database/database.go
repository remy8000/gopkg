package database

import (
	"fmt"
	"log"
	"os"
	"time"

	mysql "gorm.io/driver/mysql"
	"gorm.io/gorm"
	gormLogger "gorm.io/gorm/logger"
)

// Config for the database.
type Config struct {
	Host         string
	Port         int
	Database     string
	Username     string
	Password     string
	LoggerConfig gormLogger.Config
	TimeOut      int
}

// FKey is a foreign key descriptor.
type FKey struct {
	Model interface{}
	Args  [4]string
}

var db *gorm.DB

// DB returns the current database connection
func DB() *gorm.DB {
	return db
}

// Setup database (connection, migrations)
func Setup(config Config, tables []interface{}, fks ...FKey) error {
	uri := fmt.Sprintf(
		"%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local&timeout=%vs",
		config.Username, config.Password, config.Host, config.Port, config.Database, config.TimeOut,
	)

	newLogger := gormLogger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		config.LoggerConfig,
	)

	// be careful not ':=' but '=' in order to get a global variable
	var err error
	db, err = gorm.Open(mysql.Open(uri), &gorm.Config{
		Logger: newLogger,
	})
	if err != nil {
		return err
	}

	// configuration
	sqlDB, err := db.DB()
	if err != nil {
		return err
	}

	sqlDB.SetConnMaxIdleTime(time.Minute * 5)

	err = db.AutoMigrate(tables...)
	if err != nil {
		return err
	}

	return nil
}

// Close current database connection
func Close() error {
	sqlDB, err := db.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}
