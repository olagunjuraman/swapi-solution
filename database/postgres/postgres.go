package postgres

import (
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"time"
)

type Config struct {
	Host     string
	Port     string
	Password string
	User     string
	DBName   string
	DBURL    string
	SSLMode  string
}

func SetupDatabase(db *gorm.DB, models ...interface{}) error {
	err := db.AutoMigrate(models...).Error
	if err != nil {
		return err
	}
	return nil
}

func New(config *Config) (*gorm.DB, error) {
	var (
		db  *gorm.DB
		err error
	)
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s", config.Host, config.User, config.Password, config.DBName, config.Port, config.SSLMode)
	db, err = gorm.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}

	db.DB().SetMaxIdleConns(0)
	db.DB().SetMaxOpenConns(100)
	db.DB().SetConnMaxLifetime(time.Second * 10)

	return db, nil
}
