package database

import (
	"sync"

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
