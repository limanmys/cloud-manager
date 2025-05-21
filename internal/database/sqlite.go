package database

import (
	"log"

	_ "github.com/mattn/go-sqlite3"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func initializeSQLite() *gorm.DB {
	connection, err := gorm.Open(sqlite.Open("/tmp/testing.db"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		log.Println(err.Error())
		return nil
	}
	return connection
}
