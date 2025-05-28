package clouds

import (
	"github.com/gofiber/fiber/v2"
	"github.com/limanmys/cloud-manager-server/app/entities"
	"github.com/limanmys/cloud-manager-server/internal/database"
	"github.com/limanmys/cloud-manager-server/internal/paginator"
	"github.com/limanmys/cloud-manager-server/internal/search"
)

func Index(c *fiber.Ctx) error {

	var clouds []entities.Cloud

	db := database.Connection().Model(&entities.Cloud{}).
		Preload("Machines")

	if c.Query("search") != "" {
		search.Search(c.Query("search"), db)
	}
	page, err := paginator.New(db, c).Paginate(&clouds)
	if err != nil {
		return err
	}
	return c.JSON(page)
}
