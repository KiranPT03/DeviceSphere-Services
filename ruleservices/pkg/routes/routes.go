package routes

import (
	"github.com/gofiber/fiber/v2"

	config "github.com/kiranpt03/factorysphere/devicesphere/services/ruleservices/pkg/config"
	devices "github.com/kiranpt03/factorysphere/devicesphere/services/ruleservices/pkg/controllers/rules"
	middleware "github.com/kiranpt03/factorysphere/devicesphere/services/ruleservices/pkg/utils/middleware"
)

func Init(app *fiber.App, config *config.Config) {
	v1 := app.Group("/api/:version/rule-service", middleware.Versioning())
	{
		ruleController := devices.NewRuleController(config)
		v1.Get("/rules", ruleController.GetAllRules)
		v1.Get("/rules/:id", ruleController.GetRuleByID)
		v1.Post("/rules", ruleController.CreateRule)
		// v1.Put("/rules/:id", ruleController.UpdateDevice)
		v1.Delete("/rules/:id", ruleController.DeleteRule)
	}
}
