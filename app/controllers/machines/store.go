package machines

import (
	"github.com/gofiber/fiber/v2"
	"github.com/limanmys/cloud-manager-server/app/entities"
	"github.com/limanmys/cloud-manager-server/internal/database"
	"github.com/limanmys/cloud-manager-server/internal/validation"
	"github.com/limanmys/cloud-manager-server/pkg/random"

	"gorm.io/gorm/clause"
)

func Store(c *fiber.Ctx) error {
	var req entities.Machine
	var machine entities.Machine

	if err := c.BodyParser(&req); err != nil {
		return err
	}

	/*	old_hash, tkn, err := token.CheckExpired(c, machine.DeviceId)
		if err != nil {
			return err
		}
	*/
	database.Connection().Model(&machine).First(&machine, "device_id = ?", req.DeviceId)

	machine.Hostname = req.Hostname
	machine.DeviceId = req.DeviceId
	machine.Version = req.Version
	machine.Ostype = req.Ostype
	machine.OsName = req.OsName
	machine.OsVersion = req.OsVersion
	machine.Env = req.Env

	if req.IpAddr != "" {
		machine.IpAddr = req.IpAddr
	}

	if req.TriggerPort != 0 {
		machine.TriggerPort = req.TriggerPort
	}

	machine.Domain = req.Domain

	if !machine.IsApproved {
		machine.IsApproved = true

	}
	err := validation.Validate(&machine)
	if err != nil {
		return err
	}

	/*	if err := license.CheckMachineLicense(&machine, true); err != nil {
			log.Println(err.Error())
		}

		err = token.Validate(c, machine.DeviceId)
		if err != nil {
			return err
		}*/

	err = database.Connection().Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "device_id"}},
		UpdateAll: true,
	}).Create(&machine).Error
	if err != nil {
		return err
	}

	if machine.TriggerToken == "" {
		machine.TriggerToken = random.RandString(16)
		database.Connection().Model(&machine).Where("id = ?", machine.ID).Update("trigger_token", machine.TriggerToken)
	}

	return nil

}
