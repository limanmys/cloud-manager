package clients

import (
	"errors"
	"fmt"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/limanmys/cloud-manager-server/internal/constants"
)

func Show(c *fiber.Ctx) error {
	req := struct {
		Os   string `json:"os"`
		Arch string `json:"arch"`
	}{}
	if err := c.BodyParser(&req); err != nil {
		if c.Query("os") != "" && c.Query("arch") != "" {
			req.Os = c.Query("os")
			req.Arch = c.Query("arch")
		} else {
			return err
		}
	}
	file := fmt.Sprintf("%s/cloud-manager-client-%s.%s", constants.CLIENTS_PATH, convertArch(req.Arch), convertOs(req.Os))
	if _, err := os.Stat(file); errors.Is(err, os.ErrNotExist) {
		return fiber.NewError(fiber.StatusNotFound, err.Error())
	}
	return c.Download(file)
}
