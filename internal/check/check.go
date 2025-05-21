package check

import (
	"net/url"
	"time"

	"github.com/InVisionApp/go-health/v2"
	"github.com/InVisionApp/go-health/v2/checkers"
)

func Health() *health.Health {
	// Create a new health instance
	h := health.New()

	// Create a checker
	localUrl, _ := url.Parse("http://localhost:7878/healthcheck")
	checker, _ := checkers.NewHTTP(&checkers.HTTPConfig{
		URL: localUrl,
	})

	h.AddChecks([]*health.Config{
		{
			Name:     "cloud-manager",
			Checker:  checker,
			Interval: time.Duration(2) * time.Second,
			Fatal:    true,
		},
	})

	return h
}
