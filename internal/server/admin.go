package server

import (
	"encoding/json"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/limanmys/cloud-manager/app/routes"
	_ "github.com/limanmys/cloud-manager/internal/migrations"
)

var adminConfig = fiber.Config{
	Prefork:      false,
	BodyLimit:    1024 * 1024 * 1024,
	JSONEncoder:  json.Marshal,
	JSONDecoder:  json.Unmarshal,
	ErrorHandler: ErrorHandler,
}

func RunAdmin(test_run bool) {
	if !fiber.IsChild() {
		// Init license
	}

	app := fiber.New(adminConfig)
	app.Use(recover.New(recover.Config{EnableStackTrace: true}))
	app.Use(compress.New())
	app.Use(logger.New())

	// Admin routes
	routes.Admin(app)

	if test_run {
		log.Fatal(app.Listen("0.0.0.0:8210"))
	} else {
		listener, err := Listener()
		if err != nil {
			panic(err)
		}
		log.Fatal(app.Listener(listener))
	}
}
