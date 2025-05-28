package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/limanmys/cloud-manager-server/app/controllers/clients"
	"github.com/limanmys/cloud-manager-server/app/controllers/clouds"
	"github.com/limanmys/cloud-manager-server/app/controllers/machines"
)

func Client(app *fiber.App) {

	app.Get("/clients", clients.Index)
	app.Get("/client", clients.Show)
	app.Post("/client", clients.Show)
	app.Get("/device_id_version", clients.DeviceIdVersion)

	machineGroup := app.Group("/machines")
	{
		machineGroup.Post("/", machines.Store)
		machineGroup.Get("/:device", machines.Show)
		machineGroup.Get("/:device/env", machines.GetEnv)
		machineGroup.Get("/:device/notify_stop", machines.NotifyClientStop)
	}

	cloudGroup := app.Group("/clouds")
	{
		cloudGroup.Post("/", clouds.Store)

	}

}
