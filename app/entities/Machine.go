package entities

import (
	"fmt"

	gormjsonb "github.com/dariubs/gorm-jsonb"
	"github.com/gofiber/fiber/v2"
	"github.com/limanmys/cloud-manager-server/internal/database"
)

type Machine struct {
	Base
	Hostname              string `json:"hostname" gorm:"index" validate:"required"`
	DeviceId              string `json:"device_id" gorm:"uniqueIndex" validate:"required"`
	Version               string `json:"version" validate:"required"`
	Ostype                string `json:"ostype"`
	OsName                string `json:"os_name"`
	OsVersion             string `json:"os_version"`
	Domain                string `json:"domain"`
	IpAddr                string `json:"ip_addr"`
	TriggerPort           int    `json:"trigger_port"`
	Licensed              bool   `json:"licensed" gorm:"default:false"`
	JwtToken              string `json:"jwt_token"`
	Supplier              string `json:"supplier"`
	Description           string `json:"description"`
	MaintenanceExpiryTime int64  `json:"maintenance_expiry_time"`
	Assignee              string `json:"assignee"`
	IsApproved            bool   `json:"is_approved" gorm:"default:true"`
	Online                bool   `json:"online" gorm:"default:false"`
	TriggerToken          string `json:"trigger_token"`
	//OperatorID            *uuid.UUID `json:"operator_id"`
	//Operator              *Operator              `json:"operator"`
	Env             gormjsonb.JSONB `json:"env" gorm:"type:jsonb"`
	EnvSet          gormjsonb.JSONB `json:"env_set" gorm:"type:jsonb"`
	ApplicationType JSONB           `json:"application_type" gorm:"type:jsonb"`
	Tags            JSONB           `json:"tags" gorm:"type:jsonb"`
	Active          *bool           `json:"active" gorm:"default:true"`
}

func GetLicensedMachineCount() int64 {
	var count int64

	database.Connection().Model(&Machine{}).Where("licensed = ?", true).Count(&count)
	return count
}

func ResetLicensedCount() error {
	return database.Connection().Model(&Machine{}).
		Where("licensed = ?", true).
		Update("licensed", false).Error
}

func SetMachineLicensed(machine *Machine, licensed bool) error {

	if !machine.Online {
		return nil
	}

	return database.Connection().Model(&machine).
		Where("id = ?", machine.ID).
		Update("licensed", licensed).Error

}

func GetMachineLicensed(machine_id string) bool {
	var machine Machine
	database.Connection().Model(&machine).
		Where("device_id = ?", machine_id).
		First(&machine)

	return machine.Licensed

}

func New() fiber.Handler {
	return SetMachineStatusOnline
}

func SetMachineStatusOnline(c *fiber.Ctx) error {

	device_id := ""
	var req map[string]interface{}
	if c.Params("device", "") != "" {
		device_id = c.Params("device")
	} else {

		if err := c.BodyParser(&req); err != nil {
			return c.Next()
		}

		if req["device_id"] != nil && req["device_id"].(string) != "" {
			device_id = req["device_id"].(string)
		}
	}

	if device_id != "" {
		_ = database.Connection().Model(&Machine{}).
			Where("device_id = ?", device_id).
			Updates(map[string]interface{}{"online": true, "active": true}).Error

	}
	return c.Next()
}

func SetMachineStatusOffline(device_id string) error {

	if device_id != "" {

		_ = database.Connection().
			Exec(fmt.Sprintf("UPDATE machines set online = false where device_id = '%s' and online is true", device_id)).
			Error

		_ = database.Connection().
			Exec(fmt.Sprintf("UPDATE un_approved_machines set online = false where device_id = '%s' and online is true", device_id)).
			Error
	}
	return nil
}

func StatusController() {

	_ = database.Connection().
		Exec("UPDATE machines set online = false where updated_at < now() - interval '2 hours' and online is true").
		Error

	_ = database.Connection().
		Exec("UPDATE un_approved_machines set online = false where updated_at < now() - interval '2 hours' and online is true").
		Error

}

func DetectDuplicateMachines() {
	var machines []Machine
	/*
	   	database.Connection().Exec(`SELECT  hostname,ip_addr, count(*)
	   FROM machines
	   WHERE deleted_at is NULL
	   GROUP BY hostname,ip_addr
	   HAVING count(*) > 1`).
	   		Find(&machines)
	*/
	database.Connection().Model(&Machine{}).
		Where("deleted_at is NULL").
		Select("hostname,ip_addr, count(*)").
		Having("count(*) > 1").
		Group("hostname,ip_addr").
		Find(&machines)

	if len(machines) == 0 {
		return
	}

	for _, machine := range machines {
		var duplicates []Machine
		database.Connection().Model(&Machine{}).
			Where("ip_addr = ? and hostname = ?", machine.IpAddr, machine.Hostname).
			Order("updated_at desc").
			Find(&duplicates)
		if len(duplicates) < 2 {
			continue
		}

		for i := range duplicates {
			if i == 0 {
				continue
			}
			DeleteMachine(duplicates[i])
		}
	}
}

func DeleteMachine(machine Machine) error {

	var err error

	var tables []string
	if err := database.Connection().Table("information_schema.tables").Where("table_schema = ?", "public").Pluck("table_name", &tables).Error; err != nil {
		return err
	}

	for _, table := range tables {
		res := database.Connection().Exec(fmt.Sprintf("SELECT column_name FROM information_schema.columns WHERE table_name='%s' and column_name='machine_id';", table))

		if res.RowsAffected > 0 {
			database.Connection().Exec(fmt.Sprintf("delete from %s  where machine_id = '%s'", table, machine.ID))
		}
	}

	err = database.Connection().Delete(&Machine{}, "id = ?", machine.ID).Error
	if err != nil {
		return err
	}
	return nil
}
