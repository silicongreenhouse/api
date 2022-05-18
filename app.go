package main

import (
	"log"
	"os"

	"github.com/silicongreenhouse/api/src/executors"
	"github.com/silicongreenhouse/api/src/stores"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/joho/godotenv"

	"github.com/silicongreenhouse/api/src/sensors"
)

var App *fiber.App
var config stores.ConfigStore

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

	App.Mount("/api/sensors", sensors.Router)
	App.Mount("/api/executors", executors.Router)
}
