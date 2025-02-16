package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/recover"

	config "github.com/kiranpt03/factorysphere/devicesphere/services/deviceservices/pkg/config"
	routes "github.com/kiranpt03/factorysphere/devicesphere/services/deviceservices/pkg/routes"
	log "github.com/kiranpt03/factorysphere/devicesphere/services/deviceservices/pkg/utils/loggers"
)

func main() {
	log.Info("Starting the application.....")
	cfg, err := config.LoadConfig("./pkg/config/application.yaml")
	if err != nil {
		log.Fatal("Error while loading config, reason: %s", err.Error())
	}

	log.Info("Server Port: %s", cfg.Server.Port)

	app := fiber.New()
	app.Use(recover.New())
	app.Use(cors.New())

	routes.Init(app, cfg)

	log.Fatal("Running server on: %s", app.Listen(":"+cfg.Server.Port))

}
