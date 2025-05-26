package machines

import (
	"github.com/gofiber/fiber/v2"

	"github.com/limanmys/cloud-manager-server/app/entities"
	"github.com/limanmys/cloud-manager-server/internal/database"
	"github.com/limanmys/cloud-manager-server/internal/paginator"
	"github.com/limanmys/cloud-manager-server/internal/search"
)

func Index(c *fiber.Ctx) error {

	var machines []entities.Machine

	db := database.Connection().Model(&entities.Machine{})
	//	Preload("Operator").
	//Preload("Clouds")

	if c.Query("search") != "" {
		search.Search(c.Query("search"), db)
	}
	page, err := paginator.New(db, c).Paginate(&machines)
	if err != nil {
		return err
	}
	return c.JSON(page)
}
