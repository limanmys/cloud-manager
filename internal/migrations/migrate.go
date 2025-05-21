package migrations

import (
	"github.com/limanmys/cloud-manager/app/entities"
	"github.com/limanmys/cloud-manager/internal/database"
)

func init() {
	database.Connection().AutoMigrate(&entities.Book{})
}
