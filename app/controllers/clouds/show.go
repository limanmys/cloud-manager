package clouds

import (
	"github.com/gofiber/fiber/v2"
	"github.com/limanmys/cloud-manager-server/app/entities"
	"github.com/limanmys/cloud-manager-server/internal/database"
)

func Show(c *fiber.Ctx) error {
	var cloud entities.Cloud
	err := database.Connection().Preload("Machines").
		First(&cloud, "id = ?", c.Params("cloud")).Error

	if err != nil {
		return fiber.NewError(fiber.StatusNotFound, err.Error())
	}

	return c.JSON(cloud)
}
