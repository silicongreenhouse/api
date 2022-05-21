package main

import (
	"fmt"
	"log"
	"os"

	"github.com/silicongreenhouse/api/src/executors"
	"github.com/silicongreenhouse/api/src/stores"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/websocket/v2"
	"github.com/joho/godotenv"

	"github.com/silicongreenhouse/api/src/sensors"
)

var App *fiber.App
var config stores.ConfigStore
var socketsChannel = make(chan []byte)

func init() {
	godotenv.Load()

	config = stores.UseConfig()
	err := config.Load(os.Getenv("CONFIG_PATH"))
	if err != nil {
		log.Fatal(err)
	}

	App = fiber.New()
	App.Use(cors.New(cors.Config{
		AllowHeaders: "Content-Type, Authorization, Origin, x-access-token, XSRF-TOKEN",
	}))
	App.Use(logger.New())

	App.Mount("/api/sensors", sensors.Router)
	App.Mount("/api/executors", executors.Router)

	// Websockets requests
	App.Get("/ws_raspberry", websocket.New(func(c *websocket.Conn) {
		for {
			messageType, message, err := c.ReadMessage()
			log.Println("Message type:", messageType)
			if err != nil {
				break
			}

			log.Printf("Message: %s", message)
			go func() {
				socketsChannel <- message
			}()
			returnMessage := fmt.Sprintf("Message from server: %s", message)

			err = c.WriteMessage(messageType, []byte(returnMessage))
			if err != nil {
				break
			}
		}
		defer c.Close()
	}))

	App.Get("/ws_client", websocket.New(func(c *websocket.Conn) {
		for message := range socketsChannel {
			err = c.WriteMessage(websocket.TextMessage, []byte(message))
			if err != nil {
				log.Println(err)
				break
			}
		}
		defer c.Close()
	}))
}
