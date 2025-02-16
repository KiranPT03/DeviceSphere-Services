package routes

import (
	"github.com/gofiber/fiber/v2"

	config "github.com/kiranpt03/factorysphere/devicesphere/services/deviceservices/pkg/config"
	"github.com/kiranpt03/factorysphere/devicesphere/services/deviceservices/pkg/controllers/devices"
	"github.com/kiranpt03/factorysphere/devicesphere/services/deviceservices/pkg/utils/middleware"
)

func Init(app *fiber.App, config *config.Config) {
	v1 := app.Group("/api/:version/device-service", middleware.Versioning())
	{
		userController := devices.NewDeviceController(config)
		v1.Get("/devices", userController.GetAllDevices)
		v1.Get("/devices/:id", userController.GetDeviceByID)
		v1.Post("/devices", userController.CreateDevice)
		v1.Put("/devices/:id", userController.UpdateDevice)
		v1.Delete("/devices/:id", userController.DeleteDevice)
	}
}
