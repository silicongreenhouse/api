package sensors

import (
	"os"
	"regexp"

	"github.com/gofiber/fiber/v2"
	"github.com/silicongreenhouse/api/src/models"
)

func getSensors(c *fiber.Ctx) error {
	return c.JSON(config.State.Sensors)
}

func getSensorById(c *fiber.Ctx) error {
	id := c.Params("id")

	err, sensor := findSensor(id)
	if err != nil {
		return c.Status(400).JSON(fiber.Map {
			"err": "Sensor not found",
		})
	}

	return c.JSON(sensor)
}

func getEvents(c *fiber.Ctx) error {
	_, sensor := findSensor(c.Params("id"))

	return c.JSON(sensor.Events)
}

func getEventById(c *fiber.Ctx) error {
	_, sensor := findSensor(c.Params("id"))
	events := sensor.Events
	var desiredEvent models.Event
	found := false

	for _, event := range events {
		if event.Id == c.Params("eventId") {
			desiredEvent = event
			found = true
		}
	}

	if !found {
		return c.Status(400).JSON(fiber.Map{
			"err": "Event not found",
		})
	}

	return c.JSON(desiredEvent)
}




func editEvent(c *fiber.Ctx) error {
	returnMessage := ""
	// Getting event
	event := models.Event{}
	c.BodyParser(&event)
	sensor := models.Sensor{}

	// If id is '0' or empty create a new event
	if event.Id == "" || event.Id == "0" {
		newId, err := createNewId(c.Params("id"), event.Executor)
		if err != nil {
			c.Status(500).JSON(fiber.Map{
				"err": "Error creating event",
			})
		}

		event.Id = newId

		for index, ev := range config.State.Sensors {
			if ev.Id == c.Params("id") {
				sensor = config.State.Sensors[index]
				sensor.Events = append(sensor.Events, event)
				config.State.Sensors[index] = sensor
				break
			}
		}

		err = config.Write(os.Getenv("CONFIG_PATH"))
		if err != nil {
			return c.Status(500).JSON(fiber.Map{
				"err": "Error creating new event",
			})
		}

		returnMessage = "Event created succesfully"
		// If an id is passed then update the passed event
	} else {
		sensors := config.State.Sensors

		// Getting sensor
		for _, sen := range sensors {
			if sen.Id == c.Params("id") {
				sensor = sen
			}
		}

		// Searching events and checking if event exist
		eventFound := false
		for index, ev := range sensor.Events {
			if ev.Id == event.Id {
				regex := regexp.MustCompile(`^\D\d+`)
				executorId := string(regex.Find([]byte(event.Id)))

				if executorId != event.Executor {
					return c.Status(400).JSON(fiber.Map{
						"err": "Cannot change executor",
					})
				}

				sensor.Events[index] = event
				eventFound = true
				break
			}
		}
		if !eventFound {
			return c.Status(400).JSON(fiber.Map{
				"err": "Event not found",
			})
		}

		// updating event
		for index, sen := range config.State.Sensors {
			if sen.Id == c.Params("id") {
				config.State.Sensors[index] = sensor
			}
		}

		err := config.Write(os.Getenv("CONFIG_PATH"))
		if err != nil {
			return c.Status(400).JSON(fiber.Map{
				"err": "Error saving event",
			})
		}

		returnMessage = "Event updated succesfully"
	}

	return c.JSON(fiber.Map{
		"msg": returnMessage,
	})
}

func deleteEvent(c *fiber.Ctx) error {
	sensor := models.Sensor{}
	sensorIndex := 0
	eventFound := false

	// Getting sensor
	for i, sen := range config.State.Sensors {
		if sen.Id == c.Params("id") {
			sensor = sen
			sensorIndex = i
		}
	}

	// Finging and deleting event
	for i, ev := range sensor.Events {
		if ev.Id == c.Params("eventId") {
			eventFound = true
			deleteElement(&sensor.Events, i)
		}
	}

	if !eventFound {
		return c.Status(400).JSON(fiber.Map {
			"err": "Event not found",
		})
	}
	
	config.State.Sensors[sensorIndex] = sensor
	config.Write(os.Getenv("CONFIG_PATH"))
	
	return c.JSON(fiber.Map {
		"err": "event deleted succesfully",
	})
}
