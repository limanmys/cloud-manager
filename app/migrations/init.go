package migrations

import (
	"log"

	"github.com/limanmys/cloud-manager-server/app/entities"
	"github.com/limanmys/cloud-manager-server/internal/database"
)

func Init() {

	if err := database.Connection().AutoMigrate(&entities.Machine{}); err != nil {
		log.Fatalln("error when making migrations, ", err.Error())
	}

}
