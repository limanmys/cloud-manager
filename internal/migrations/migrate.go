package migrations

import (
	"github.com/limanmys/fiber-app-template/app/entities"
	"github.com/limanmys/fiber-app-template/internal/database"
)

func init() {
	database.Connection().AutoMigrate(&entities.Book{})
}
