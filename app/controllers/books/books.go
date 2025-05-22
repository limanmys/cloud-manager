package books

import (
	"errors"

	"github.com/google/uuid"
	"github.com/limanmys/cloud-manager-server/app/entities"
	"github.com/limanmys/cloud-manager-server/internal/database"
	"github.com/limanmys/cloud-manager-server/internal/paginator"
	"github.com/limanmys/cloud-manager-server/internal/search"
	"github.com/limanmys/cloud-manager-server/internal/validation"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm/clause"
)

func Create(c *fiber.Ctx) error {
	// Parse request body
	item := entities.Book{}
	if err := c.BodyParser(&item); err != nil {
		return err
	}
	// Validate request
	err := validation.Validate(&item)
	if err != nil {
		return err
	}
	// Create item on database
	err = database.Connection().Model(&entities.Book{}).Create(&item).Error
	if err != nil {
		return err
	}
	return c.JSON(item)
}

func Index(c *fiber.Ctx) error {
	// Create empty slice
	items := []entities.Book{}
	// Get all items from database
	db := database.Connection().Model(&entities.Book{}).Preload(clause.Associations)
	// If search query exists
	if c.Query("search") != "" {
		search.Search(c.Query("search"), db)
	}
	// Paginate return value
	page, err := paginator.New(db, c).Paginate(&items)
	if err != nil {
		return err
	}
	return c.JSON(page)
}

func Show(c *fiber.Ctx) error {
	// Check is id exists
	if len(c.Params("id")) <= 0 {
		return errors.New("please set id as parameter")
	}
	// Get item with id from parameters
	item := entities.Book{}
	err := database.Connection().Model(&item).Preload(clause.Associations).Where("id = ?", c.Params("id")).First(&item).Error
	if err != nil {
		return err
	}

	return c.JSON(item)
}

func Update(c *fiber.Ctx) error {
	// Check is id exists
	if len(c.Params("id")) <= 0 {
		return errors.New("please set id as parameter")
	}
	// Parse request body
	request := entities.Book{}
	if err := c.BodyParser(&request); err != nil {
		return err
	}
	// Get item with id
	item := entities.Book{}
	err := database.Connection().Model(&item).Where("id = ?", c.Params("id")).First(&item).Error
	if err != nil {
		return err
	}
	// Update local item
	err = database.Connection().Model(&item).Updates(&request).Error
	if err != nil {
		return err
	}

	return c.JSON(request)
}

func Delete(c *fiber.Ctx) error {
	// Check is id exists
	if len(c.Params("id")) <= 0 {
		return errors.New("please set id as parameter")
	}
	// Parse id for validity check
	_, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return err
	}
	// Delete item
	err = database.Connection().Model(&entities.Book{}).Delete("id = ?", c.Params("id")).Error
	if err != nil {
		return err
	}

	return c.JSON("Item deleted successfully.")
}
