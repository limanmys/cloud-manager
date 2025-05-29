package clouds

import (
	"github.com/gofiber/fiber/v2"
	"github.com/limanmys/cloud-manager-server/app/entities"
	"github.com/limanmys/cloud-manager-server/internal/database"
	"gorm.io/gorm/clause"
)

type cloud_info struct {
	DeviceId     string                 `json:"device_id"`
	RegisterInfo map[string]interface{} `json:"register_info"`
}

func Store(c *fiber.Ctx) error {
	var req cloud_info
	var cloud entities.Cloud
	var machines []entities.Machine
	var machine_info entities.Machine
	var hosts []string

	if err := c.BodyParser(&req); err != nil {
		return err
	}

	if _, ok := req.RegisterInfo["hosts"].([]interface{}); ok {

		var str_list []string
		for _, item := range req.RegisterInfo["hosts"].([]interface{}) {
			str_list = append(str_list, item.(string))
		}

		if len(str_list) == 1 {
			cloud.Name = str_list[0]
		}
		hosts = str_list
	}
	if len(hosts) > 0 {
		database.Connection().Model(&entities.Machine{}).Find(&machines, "hostname in (?) and cloud_type = ?", hosts, req.RegisterInfo["type"].(string))

	}

	err := database.Connection().Model(&machine_info).First(&machine_info, "device_id = ?", req.DeviceId).Error
	if err != nil {
		return fiber.NewError(fiber.StatusNotFound, "machine not found:%s", err.Error())
	}

	if _, ok := req.RegisterInfo["type"].(string); !ok {
		return fiber.NewError(fiber.StatusNotAcceptable, "invalid cloud type")
	}

	cloud.Name = machine_info.Hostname

	res := database.Connection().Model(&cloud).First(&cloud, "name = ?", cloud.Name)

	if res.RowsAffected == 0 {

		for _, machine := range machines {
			var items []map[string]string
			database.Connection().Table("cloud_machines").Where("machine_id = ?", machine.ID).
				Find(&items)

			if len(items) > 0 {
				var cloud_info entities.Cloud
				database.Connection().Model(&cloud_info).First(&cloud_info, "id = ?", items[0]["cloud_id"])
				cloud.Name = cloud_info.Name

			}
		}

	}

	cloud.Type = req.RegisterInfo["type"].(string)
	if cloud.Name == "" {
		return fiber.NewError(fiber.StatusNotFound, "cloud not found")
	}

	err = database.Connection().Clauses(clause.OnConflict{
		UpdateAll: true,
	}).Create(&cloud).Error
	if err != nil {
		return err
	}

	if len(hosts) < 1 {
		return fiber.NewError(fiber.StatusOK, "cloud register completed")
	}

	for _, machine := range machines {
		database.Connection().Model(&cloud).Debug().Association("Machines").Append(&machine)
	}

	return fiber.NewError(fiber.StatusOK, "cloud register completed")

}
