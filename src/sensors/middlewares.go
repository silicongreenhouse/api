package sensors

import (
	"github.com/gofiber/fiber/v2"
	"github.com/silicongreenhouse/api/src/models"
)

func checkSensor(c *fiber.Ctx) error {
	id := c.Params("id")

	err, _ := findSensor(id)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"err": "Sensor not found",
		})
	}

	return c.Next()
}

func validateEventData(c *fiber.Ctx) error {
	event := models.Event{}
	executors := config.Executors

	err := c.BodyParser(&event)

	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"err": "Invalid values types",
		})
	}

	// Check if executor exist
	executorFound := false
	for _, executor := range executors {
		if executor.Id == event.Executor {
			executorFound = true
			break
		}
	}

	if !executorFound {
		return c.Status(400).JSON(fiber.Map{
			"err": "Executor not found",
		})
	}

	// Check if the state is correct
	if event.State != "on" && event.State != "off" {
		return c.Status(400).JSON(fiber.Map{
			"err": "Invalid state value",
		})
	}

	// check Value only if the sensor is water level
	if c.Params("id") == "s3" {
		if event.Value != "high" && event.Value != "low" {
			return c.Status(400).JSON(fiber.Map{
				"err": "Invalid value",
			})
		}
	} else {
		if event.Above == 0 && event.Equal == 0 && event.Below == 0 {
			return c.Status(400).JSON(fiber.Map{
				"err": "Have to put some condition values",
			})
		}
	}

	return c.Next()
}
