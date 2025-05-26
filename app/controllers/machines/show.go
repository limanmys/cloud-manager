package machines

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/limanmys/cloud-manager-server/app/entities"
	"github.com/limanmys/cloud-manager-server/internal/database"
)

func Show(c *fiber.Ctx) error {
	var machine entities.Machine
	connection := database.Connection().Preload("Clouds")
	if c.Params("device", "") != "" {
		connection = connection.First(&machine, "device_id = ?", c.Params("device"))
	} else if c.Params("machine", "") != "" {
		connection = connection.First(&machine, "id = ?", c.Params("machine"))
	}
	if connection.Error != nil {
		return fiber.NewError(fiber.StatusNotFound, connection.Error.Error())
	}
	if machine.MaintenanceExpiryTime == 0 {
		machine.MaintenanceExpiryTime = time.Now().UnixMilli() + time.Hour.Milliseconds()*24*365
	}
	return c.JSON(machine)
}
