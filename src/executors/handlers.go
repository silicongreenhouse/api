package executors

import (
	"github.com/gofiber/fiber/v2"
)

func getExecutors(c *fiber.Ctx) error {
	return c.JSON(config.State.Executors)
}


