package machines

import (
	"github.com/gofiber/fiber/v2"
	"github.com/limanmys/cloud-manager-server/app/entities"
	"github.com/limanmys/cloud-manager-server/internal/database"
)

func GetEnv(c *fiber.Ctx) error {

	var err error
	var machine entities.Machine
	if c.Params("device", "") != "" {
		err = database.Connection().First(&machine, "device_id = ?", c.Params("device")).Error
	} else if c.Params("machine", "") != "" {
		err = database.Connection().First(&machine, "id = ?", c.Params("machine")).Error
	}
	if err != nil {
		return fiber.NewError(fiber.StatusNotFound, err.Error())
	}

	if c.Params("device", "") != "" {
		return c.JSON(machine.EnvSet)
	}
	return c.JSON(machine.Env)
}
