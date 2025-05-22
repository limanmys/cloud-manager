package server

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/limanmys/cloud-manager-server/app/routes"
	"github.com/limanmys/cloud-manager-server/internal/socket"
)

var clientConfig = fiber.Config{
	Prefork:      false,
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

	app.Use("/ws", func(c *fiber.Ctx) error {
		if websocket.IsWebSocketUpgrade(c) {
			c.Locals("allowed", true)
			//todo: check authentication here
			c.Locals("machine_id", string(c.Request().Header.Peek("machine_id")))
			return c.Next()
		}
		return fiber.ErrUpgradeRequired
	})

	go socket.WebSocketHandler()

	socket.Init(app)

	// Set client routes
	routes.Client(app)

	// Set healthcheck endpoint
	/*h := check.Health()
	go h.Start()
	app.Get("/healthcheck", adaptor.HTTPHandler(handlers.NewJSONHandlerFunc(h, nil)))
	*/
	go func() {
		for {
			time.Sleep(2 * time.Second)
			res, err := socket.Send("1111", "test", "data")
			if err != nil {
				fmt.Println("error:", err.Error())
				continue
			}
			fmt.Println("res", res)
		}

	}()
	log.Fatal(app.ListenTLS(fmt.Sprintf("%s:%d", "0.0.0.0", 8211), "/opt/cloud-manager-server/keys/cloud-manager-server.pem", "/opt/cloud-manager-server/keys/cloud-manager-server.key"))
}

func handler(f http.HandlerFunc) http.Handler {
	return http.HandlerFunc(f)
}
