package machines

import (
	"github.com/gofiber/fiber/v2"
	"github.com/limanmys/cloud-manager-server/app/entities"
	"github.com/limanmys/cloud-manager-server/internal/database"
)

func NotifyClientStop(c *fiber.Ctx) error {
	err := database.Connection().Model(&entities.Machine{}).Where("device_id = ?", c.Params("device")).
		Update("online", "false").Error
	if err != nil {
		return err
	}

	return c.JSON("Status updated successfully.")
}
