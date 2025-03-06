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
		deviceController := devices.NewDeviceController(config)
		deviceDiscoveryController := devices.NewDeviceDiscoveryController(config)
		v1.Get("/devices", deviceController.GetAllDevices)
		v1.Get("/devices/:id", deviceController.GetDeviceByID)
		v1.Post("/devices", deviceController.CreateDevice)
		v1.Put("/devices/:id", deviceController.UpdateDevice)
		v1.Delete("/devices/:id", deviceController.DeleteDevice)
		// Properties
		v1.Get("/devices/:id/properties", deviceController.GetAllPropertiesByDeviceID)
		v1.Get("/device-discovery/kepware/devices", deviceDiscoveryController.GetKepwareDeviceData)
	}
}
