package database

import (
	"fmt"
	"os"
	"sync"

	"github.com/joho/godotenv"
	"gorm.io/gorm"
)

var once sync.Once
var connection *gorm.DB

func Connection() *gorm.DB {
	once.Do(func() {
		connection = initialize()
	})

	return connection
}

func initialize() *gorm.DB {

	e := godotenv.Load()
	if e != nil {
		fmt.Print(e)
		os.Exit(1)
	}

	/*switch os.Getenv("DB_DRIVER") {
	case "postgres":
		return initializePostgres()
	case "mysql":
		return initializeMysql()
	case "sqlite":
		return initializeSQLite()
	default:
		log.Fatalln("You must specify a database driver. Choices are 'postgres' / 'mysql' / 'sqlite'")
		return nil
	}*/

	return initializePostgres()
}
