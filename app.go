package main

import (
	"encoding/json"
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
var streamDataChannel = make(chan []byte)
var remoteControllerChannel = make(chan []byte)

var clientConnected = false
var clientControllerConnected = false
var raspberryConnected = false
var raspberryControllerConnected = false

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
		raspberryConnected = true
		for {
			_, message, err := c.ReadMessage()
			if err != nil {
				log.Println(err)
				raspberryConnected = false
				break
			}

			if clientConnected {
				streamDataChannel <- message
			}

			returnMessage, jsonError := json.Marshal(fiber.Map{
				"msg": "Data sent succesfully bitch",
			})
			if jsonError != nil {
				log.Println(jsonError)
				c.Close()
				break
			}
			err = c.WriteMessage(websocket.TextMessage, returnMessage)
			if err != nil {
				log.Println(err)
				raspberryConnected = false
				break
			}
		}

		defer c.Close()
	}))

	App.Get("/ws_raspberry_controller", websocket.New(func(c *websocket.Conn) {
		raspberryControllerConnected = true
		for message := range remoteControllerChannel {
			log.Println(message)
			err := c.WriteMessage(websocket.TextMessage, []byte(message))
			if err != nil {
				raspberryConnected = false
				break
			}
		}
		defer c.Close()
	}))

	App.Get("/ws_client", websocket.New(func(c *websocket.Conn) {
		clientConnected = true
		for message := range streamDataChannel {
			err := c.WriteMessage(websocket.TextMessage, []byte(message))
			if err != nil {
				clientConnected = false
				break
			}
		}
		defer c.Close()
	}))

	App.Get("/ws_client_controller", websocket.New(func(c *websocket.Conn) {
		clientControllerConnected = true
		for {
			_, message, err := c.ReadMessage()
			if err != nil {
				clientControllerConnected = false
				break
			}

			if raspberryControllerConnected {
				remoteControllerChannel <- message
			}

			returnMessage, jsonErr := json.Marshal(fiber.Map{
				"msg": "Data sent succesfully",
			})
			if jsonErr != nil {
				c.Close()
				break
			}

			err = c.WriteMessage(websocket.TextMessage, returnMessage)
			if err != nil {
				clientConnected = false
				continue
			}
		}
		defer c.Close()
	}))
}
