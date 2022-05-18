package executors

import (
	"github.com/gofiber/fiber/v2"
	"github.com/silicongreenhouse/api/src/stores"
)

var Router *fiber.App
var config stores.ConfigStore

func init() {
	Router = fiber.New()
	Router.Get("/", getExecutors)

	config = stores.UseConfig()
}
