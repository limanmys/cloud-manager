package server

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/InVisionApp/go-health/v2/handlers"
	"github.com/gofiber/adaptor/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/limanmys/cloud-manager/app/routes"
	"github.com/limanmys/cloud-manager/internal/check"
)

var clientConfig = fiber.Config{
	Prefork:      true,
	BodyLimit:    1024 * 1024 * 1024,
	JSONEncoder:  json.Marshal,
	JSONDecoder:  json.Unmarshal,
	ErrorHandler: ErrorHandler,
}

func RunClient() {
	// Set some config
	app := fiber.New(clientConfig)
	app.Use(recover.New(recover.Config{EnableStackTrace: true}))
	app.Use(compress.New())
	app.Use(logger.New())

	// Set client routes
	routes.Client(app)

	// Set healthcheck endpoint
	h := check.Health()
	go h.Start()
	app.Get("/healthcheck", adaptor.HTTPHandler(handlers.NewJSONHandlerFunc(h, nil)))

	log.Fatal(app.ListenTLS(fmt.Sprintf("%s:%d", "0.0.0.0", 7878), "/opt/cloud-manager/keys/cloud-manager.pem", "/opt/cloud-manager/keys/cloud-manager.key"))
}

func handler(f http.HandlerFunc) http.Handler {
	return http.HandlerFunc(f)
}
