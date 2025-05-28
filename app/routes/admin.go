package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/limanmys/cloud-manager-server/app/controllers/clouds"
	"github.com/limanmys/cloud-manager-server/app/controllers/machines"
)

func Admin(app *fiber.App) {

	machinesGroup := app.Group("/machines")
	{
		machinesGroup.Get("/", machines.Index)
		machinesGroup.Get("/:machine", machines.Show)
	}

	cloudsGroup := app.Group("/clouds")
	{
		cloudsGroup.Get("/", clouds.Index)
		cloudsGroup.Get("/:cloud", clouds.Show)
	}

}
