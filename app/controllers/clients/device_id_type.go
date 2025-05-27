package clients

import (
	"github.com/gofiber/fiber/v2"
	"github.com/limanmys/cloud-manager-server/pkg/config"
)

func DeviceIdVersion(c *fiber.Ctx) error {

	res := "1.0"
	res = config.Get("DEVICE_ID_VERSION", "1.0")
	return c.JSON(map[string]string{"device_id_version": res})
}
