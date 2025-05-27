package clients

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path"
	"regexp"
	"strings"

	gormjsonb "github.com/dariubs/gorm-jsonb"
	"github.com/gofiber/fiber/v2"
	"github.com/limanmys/cloud-manager-server/internal/constants"
	"github.com/limanmys/cloud-manager-server/pkg/config"
)

func Index(c *fiber.Ctx) error {
	var client_ips []string
	hashes := make(gormjsonb.JSONB)
	clients, err := os.ReadDir(constants.CLIENTS_PATH)
	if err != nil {
		return err
	}
	allow_update := config.Get("ALLOW_UPDATE", "")
	if allow_update != "" {
		client_ips = strings.Split(allow_update, ",")
	}

	found := false
	for _, cl := range client_ips {
		if c.IP() == strings.TrimSpace(cl) {
			found = true
		}
	}

	if !found && len(client_ips) > 0 {
		return fmt.Errorf("update not allowed")
	}

	var re = regexp.MustCompile(`(?m)^.*\.sum$`)
	for _, client := range clients {
		if !client.IsDir() {
			if re.MatchString(client.Name()) {
				var h Hashes
				hash, err := os.ReadFile((path.Join(constants.CLIENTS_PATH, client.Name())))
				if err != nil {
					log.Printf("unable to read file: %v", err)
				}
				err = json.Unmarshal(hash, &h)
				if err != nil {
					log.Printf("unable to parse hashes: %v", err)
					continue
				}
				hashes[convert(client.Name())] = h
			}
		}
	}
	return c.JSON(hashes)
}
