package sensors

import (
	"github.com/gofiber/fiber/v2"
	"github.com/silicongreenhouse/api/src/stores"
)

var Router *fiber.App
var config stores.ConfigStore

func init() {
	Router = fiber.New()
	Router.Get("/:id/events", checkSensor, getEvents)
	Router.Put("/:id/events", checkSensor, validateEventData, editEvent)
	Router.Delete("/:id/events/:eventId", checkSensor, deleteEvent)
	
	config = stores.UseConfig()
}