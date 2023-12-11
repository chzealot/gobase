package database

import (
	"database/sql"
	"fmt"
	"github.com/pkg/errors"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"os"
	"strconv"
	"time"
)

type Database struct {
	DB   *gorm.DB
	Host string
	Name string
}

func NewDatabase() *Database {
	return &Database{}
}

func (d *Database) InitDbFromEnv() error {
	user, ok := os.LookupEnv("MYSQL_USERNAME")
	if !ok {
		return fmt.Errorf("InitDbFromEnv, env.MYSQL_USERNAME is empty")
	}
	pass, ok := os.LookupEnv("MYSQL_PASSWORD")
	if !ok {
		return fmt.Errorf("InitDbFromEnv, env.MYSQL_PASSWORD is empty")
	}
	host, ok := os.LookupEnv("MYSQL_HOST")
	if !ok {
		return fmt.Errorf("InitDbFromEnv, env.MYSQL_HOST is empty")
	}
	portStr := os.Getenv("MYSQL_PORT")
	if len(portStr) == 0 {
		portStr = "3306"
	}
	dbName, ok := os.LookupEnv("MYSQL_DATABASE")
	if !ok {
		return fmt.Errorf("InitDbFromEnv, env.MYSQL_DATABASE is empty")
	}
	port, err := strconv.Atoi(portStr)
	if err != nil {
		return errors.Wrap(err, "InitDbFromEnv, env.MYSQL_PORT is invalid")
	}

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		user, pass, host, port, dbName)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("open database(%s)", dbName))
	}
	var sqlDB *sql.DB
	sqlDB, err = db.DB()
	if err != nil {
		return errors.Wrap(err, "db.DB()")
	}
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	d.DB = db
	d.Host = host
	d.Name = dbName
	return nil
}
