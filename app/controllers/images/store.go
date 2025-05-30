package images

import (
	"github.com/gofiber/fiber/v2"
	"github.com/limanmys/cloud-manager-server/app/entities"
	"github.com/limanmys/cloud-manager-server/internal/database"
	"gorm.io/gorm/clause"
)

type image_data struct {
	DeviceId string           `json:"device_id"`
	Images   []entities.Image `json:"images"`
}

func Store(c *fiber.Ctx) error {
	var req image_data
	var machine_info entities.Machine
	var cloud entities.Cloud

	if err := c.BodyParser(&req); err != nil {
		return err
	}
	err := database.Connection().Model(&machine_info).First(&machine_info, "device_id = ?", req.DeviceId).Error
	if err != nil {
		return fiber.NewError(fiber.StatusNotFound, "machine not found:%s", err.Error())
	}

	var item map[string]string
	err = database.Connection().Table("cloud_machines").Where("machine_id = ?", machine_info.ID).
		First(&item).Error

	if err != nil {
		return fiber.NewError(fiber.StatusNotFound, "cloud not found:%s", err.Error())
	}

	err = database.Connection().Model(&cloud).First(&cloud, "id = ?", item["cloud_id"]).Error
	if err != nil {
		return fiber.NewError(fiber.StatusNotFound, "cloud not found:%s", err.Error())
	}

	err = database.Connection().Clauses(clause.OnConflict{
		UpdateAll: true,
	}).Create(&req.Images).Error
	if err != nil {
		return err
	}

	for _, image := range req.Images {
		database.Connection().Model(&cloud).Debug().Association("Machines").Append(&image)
	}

	return fiber.NewError(fiber.StatusOK, "cloud register completed")

}
