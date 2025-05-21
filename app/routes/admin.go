package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/limanmys/fiber-app-template/app/controllers/books"
)

func Admin(app *fiber.App) {
	group := app.Group("/books")
	{
		// Create item on database
		group.Post("/books", books.Create)
		// Get all items from database
		group.Get("/books", books.Index)
		// Show single item
		group.Get("/books/:id", books.Show)
		// Update item with id
		group.Patch("/books/:id", books.Update)
		// Delete item with id
		group.Delete("/books/:id", books.Delete)
	}
}
