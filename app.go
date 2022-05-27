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

	"github.com/silicongreenhouse/api/src/sensors"
)

var App *fiber.App
var config stores.ConfigStore
var streamDataChannel = make(chan []byte)
var remoteControllerChannel = make(chan []byte)

func init() {
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

	App.Static("/control_panel", staticFolder)

	App.Mount("/api/sensors", sensors.Router)
	App.Mount("/api/executors", executors.Router)

	// Websockets requests
	App.Get("/ws_trigger", websocket.New(func(c *websocket.Conn) {
		stores.RaspberryConnected = true
		defer c.Close()
		for {
			select {
			case msg := <-config.SignalChannel:
				if msg {
					returnMessage, jsonError := json.Marshal(fiber.Map{
						"msg": "ConfigChanged",
					})
					if err != nil {
						log.Println(jsonError)
						c.Close()
						break
					}

					err = c.WriteMessage(websocket.TextMessage, returnMessage)
					if err != nil {
						log.Println(err)
						stores.RaspberryConnected = false
						break
					}
				}
			}
		}
	}))

	App.Get("/ws_raspberry", websocket.New(func(c *websocket.Conn) {
		stores.RaspberryConnected = true
		for {
			_, message, err := c.ReadMessage()
			if err != nil {
				log.Println(err)
				stores.RaspberryConnected = false
				break
			}

			if stores.ClientConnected {
				streamDataChannel <- message
			}
		}

		defer c.Close()
	}))

	App.Get("/ws_raspberry_controller", websocket.New(func(c *websocket.Conn) {
		stores.RaspberryControllerConnected = true
		for message := range remoteControllerChannel {
			log.Println(message)
			err := c.WriteMessage(websocket.TextMessage, []byte(message))
			if err != nil {
				stores.RaspberryConnected = false
				break
			}
		}
		defer c.Close()
	}))

	App.Get("/ws_client", websocket.New(func(c *websocket.Conn) {
		stores.ClientConnected = true
		for message := range streamDataChannel {
			err := c.WriteMessage(websocket.TextMessage, []byte(message))
			if err != nil {
				stores.ClientConnected = false
				break
			}
		}
		defer c.Close()
	}))

	App.Get("/ws_client_controller", websocket.New(func(c *websocket.Conn) {
		stores.ClientControllerConnected = true
		for {
			_, message, err := c.ReadMessage()
			if err != nil {
				stores.ClientControllerConnected = false
				break
			}

			if stores.RaspberryControllerConnected {
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
				stores.ClientConnected = false
				continue
			}
		}
		defer c.Close()
	}))
}
