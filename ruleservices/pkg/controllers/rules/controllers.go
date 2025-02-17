package devices

import (
	"fmt"
	"reflect"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"

	models "github.com/kiranpt03/factorysphere/devicesphere/services/ruleservices/pkg/models"
	log "github.com/kiranpt03/factorysphere/devicesphere/services/ruleservices/pkg/utils/loggers"
	resterrors "github.com/kiranpt03/factorysphere/devicesphere/services/ruleservices/pkg/utils/resterrors"

	// Import the PostgreSQL repository package
	config "github.com/kiranpt03/factorysphere/devicesphere/services/ruleservices/pkg/config"
	postgres "github.com/kiranpt03/factorysphere/devicesphere/services/ruleservices/pkg/databases/postgres"
)

type RuleController struct {
	repository *postgres.PostgreSQLRepository
}

func NewRuleController(config *config.Config) *RuleController {
	// Create a new PostgreSQL repository instance
	repository, err := postgres.NewPostgreSQLRepository(config)
	if err != nil {
		panic(err)
	}

	return &RuleController{repository: repository}
}

func (rc *RuleController) getRuleData(ruleId string) (models.Rule, error) {
	rule := models.Rule{}
	query := `
				SELECT
					r.*,
					c.id AS condition_id,
					c.position,
					c.type,
					c.device_id,
					c.property_id,
					c.operator_id,
					c.operator_symbol,
					c.value
				FROM
					rules r
				LEFT JOIN
					conditions c ON r.id = c.rule_id
				WHERE
					r.id = $1
			`

	result, queryErr := rc.repository.ExecuteQuery(query, ruleId)
	if queryErr != nil {
		log.Error("Error while executing query: %v", queryErr)
		return models.Rule{}, error(queryErr)
	}
	log.Debug("Result: %v", result)
	if len(result) > 0 {
		rule.ID = result[0]["id"].(string)
		rule.Name = result[0]["name"].(string)
		rule.Severity = result[0]["severity"].(string)
		rule.Status = result[0]["status"].(string)
		rule.CreatedAt = result[0]["created_at"].(string)
		rule.UpdatedAt = result[0]["updated_at"].(string)

		for _, condMap := range result {
			condition := models.Condition{
				ID:             condMap["condition_id"].(string),
				Position:       condMap["position"].(string),
				Type:           condMap["type"].(string),
				DeviceId:       condMap["device_id"].(string),
				PropertyId:     condMap["property_id"].(string),
				OperatorId:     condMap["operator_id"].(string),
				OperatorSymbol: condMap["operator_symbol"].(string),
				Value:          condMap["value"].(string),
			}
			rule.Conditions = append(rule.Conditions, condition)
		}
	}

	return rule, nil
}

func (rc *RuleController) GetAllRules(c *fiber.Ctx) error {
	// Use the repository to get all rules

	ruleList := []models.Rule{}
	ruleIds, err := rc.repository.GetAll("rules")
	if err != nil {
		return resterrors.SendInternalServerError(c)
	}

	log.Debug("Rules: %v", ruleIds)
	for _, rule := range ruleIds {
		log.Debug("Rule ID: %s", rule["id"])
		ruleId := rule["id"].(string)

		rule, err := rc.getRuleData(ruleId)
		if err != nil {
			return resterrors.SendInternalServerError(c)
		}

		// rules = append(rules, rule)
		if !reflect.DeepEqual(rule, models.Rule{}) {
			ruleList = append(ruleList, rule)
		}
	}

	return c.JSON(ruleList)
}

func (rc *RuleController) GetRuleByID(c *fiber.Ctx) error {
	// Get the rule ID from the request parameters
	id := c.Params("id")

	// Use the repository to get the rule by ID
	rule, err := rc.getRuleData(id)
	if err != nil {
		return resterrors.SendNotFoundError(c, "No rule found")
	}
	if reflect.DeepEqual(rule, models.Rule{}) {
		return resterrors.SendNotFoundError(c, "No rule found")
	}

	return resterrors.SendOK(c, rule)
}

// @Summary      Create a new device
// @Description  Create a new device with the given parameters
// @Tags         devices
// @Accept       json
// @Produce      json
// @Param        device  body      models.Device  true  "Device"
// @Success      201      {object}   map[string]interface{}  "Device created"
// @Failure      400      {object}   httperrors.HTTPError  "Bad request"
// @Failure      409      {object}   httperrors.HTTPError  "Device already exists"
// @Failure      500      {object}   httperrors.HTTPError  "Internal server error"
// @Router       /devices [post]
func (rc *RuleController) CreateRule(c *fiber.Ctx) error {
	log.Info("Creating a new rule")
	var rule models.Rule
	err := c.BodyParser(&rule)
	log.Debug("Err: %e", err)
	if err != nil {
		return resterrors.SendBadRequestError(c)
	}
	// Validate the rule
	validate := validator.New()
	err = validate.Struct(rule)
	if err != nil {
		return resterrors.SendBadRequestError(c)
	}
	log.Debug("Rule: %v", rule)
	ruleMap := make(map[string]interface{})
	ruleMap["name"] = rule.Name
	ruleMap["severity"] = rule.Severity
	ruleMap["status"] = rule.Status
	ruleMap["created_at"] = rule.CreatedAt
	ruleMap["updated_at"] = rule.UpdatedAt
	ruleId, ruleErr := rc.repository.Create("rules", ruleMap)
	if ruleErr != nil {
		return resterrors.SendInternalServerError(c)
	}
	conditions := make([]map[string]interface{}, len(rule.Conditions))
	for i, condition := range rule.Conditions {
		conditionMap := make(map[string]interface{})
		conditionMap["position"] = condition.Position
		conditionMap["type"] = condition.Type
		conditionMap["device_id"] = condition.DeviceId
		conditionMap["property_id"] = condition.PropertyId
		conditionMap["operator_id"] = condition.OperatorId
		conditionMap["operator_symbol"] = condition.OperatorSymbol
		conditionMap["value"] = condition.Value
		conditionMap["rule_id"] = ruleId
		conditionId, conditionErr := rc.repository.Create("conditions", conditionMap)
		if conditionErr != nil {
			return resterrors.SendInternalServerError(c)
		}
		conditionMap["id"] = conditionId
		conditions[i] = conditionMap
	}

	ruleMap["condition"] = conditions

	fmt.Println(ruleMap)

	// Convert ruleMap to rule model
	var newRule models.Rule
	newRule.ID = ruleMap["id"].(string)
	newRule.Name = ruleMap["name"].(string)
	newRule.Severity = ruleMap["severity"].(string)
	newRule.Status = ruleMap["status"].(string)
	newRule.CreatedAt = ruleMap["created_at"].(string)
	newRule.UpdatedAt = ruleMap["updated_at"].(string)

	conditions = ruleMap["condition"].([]map[string]interface{})
	newRule.Conditions = make([]models.Condition, len(conditions))

	for i, condition := range conditions {
		newRule.Conditions[i] = models.Condition{
			ID:             condition["id"].(string),
			Position:       condition["position"].(string),
			Type:           condition["type"].(string),
			DeviceId:       condition["device_id"].(string),
			PropertyId:     condition["property_id"].(string),
			OperatorId:     condition["operator_id"].(string),
			OperatorSymbol: condition["operator_symbol"].(string),
			Value:          condition["value"].(string),
		}
	}

	return resterrors.SendCreated(c, newRule)
}

// func (dc *RuleController) UpdateDevice(c *fiber.Ctx) error {
// 	// Get the device ID from the request parameters
// 	id := c.Params("id")

// 	// Create a new device instance from the request body
// 	var device map[string]interface{}
// 	err := c.BodyParser(&device)
// 	if err != nil {
// 		return resterrors.SendBadRequestError(c)
// 	}

// 	// Use the repository to update the device
// 	affectedRows, err := dc.repository.Update("devices", id, device)
// 	if err != nil {
// 		return resterrors.SendInternalServerError(c)
// 	}

// 	if affectedRows == 0 {
// 		return resterrors.SendNotFoundError(c, "No device found")
// 	}

// 	return resterrors.SendOK(c, affectedRows)
// }

func (rc *RuleController) DeleteRule(c *fiber.Ctx) error {
	// Get the rule ID from the request parameters
	id := c.Params("id")

	// Use the repository to delete the rule
	ruleData, getErr := rc.repository.CheckExist("rules", "id", id)
	if getErr != nil {
		return resterrors.SendInternalServerError(c)
	}
	if !ruleData {
		return resterrors.SendNotFoundError(c, "Rule not found")
	}
	condErr := rc.repository.Delete("conditions", "rule_id", id)
	log.Debug("Delete Error %v", condErr)
	if condErr != nil {
		return resterrors.SendInternalServerError(c)
	}

	err := rc.repository.Delete("rules", "id", id)
	log.Debug("Delete Error %v", err)
	if err != nil {
		return resterrors.SendInternalServerError(c)
	}
	return resterrors.SendNoContent(c)
}
