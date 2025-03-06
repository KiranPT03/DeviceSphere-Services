package devices

import (
	"fmt"
	"reflect"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"

	models "github.com/kiranpt03/factorysphere/devicesphere/services/deviceservices/pkg/models"
	log "github.com/kiranpt03/factorysphere/devicesphere/services/deviceservices/pkg/utils/loggers"
	resterrors "github.com/kiranpt03/factorysphere/devicesphere/services/deviceservices/pkg/utils/resterrors"

	// Import the PostgreSQL repository package
	config "github.com/kiranpt03/factorysphere/devicesphere/services/deviceservices/pkg/config"
	postgres "github.com/kiranpt03/factorysphere/devicesphere/services/deviceservices/pkg/databases/postgres"
)

type DeviceController struct {
	repository *postgres.PostgreSQLRepository
}

func NewDeviceController(config *config.Config) *DeviceController {
	// Create a new PostgreSQL repository instance
	repository, err := postgres.NewPostgreSQLRepository(config)
	if err != nil {
		panic(err)
	}

	return &DeviceController{repository: repository}
}

func (dc *DeviceController) GetAllPropertiesByDeviceID(c *fiber.Ctx) error {
	deviceId := c.Params("id")
	query := `
                SELECT
                        id, reference_id, name, unit, state, status, data_type, value, threshold
                FROM
                        properties
                WHERE
                        device_id = $1
        `

	results, dbErr := dc.repository.ExecuteQuery(query, deviceId)
	if dbErr != nil {
		return resterrors.SendInternalServerError(c)
	}

	properties := make([]models.Property, 0, len(results))

	for _, result := range results {
		property := models.Property{
			ID:          fmt.Sprintf("%v", result["id"]),
			ReferenceID: fmt.Sprintf("%v", result["reference_id"]),
			Name:        fmt.Sprintf("%v", result["name"]),
			Unit:        fmt.Sprintf("%v", result["unit"]),
			State:       fmt.Sprintf("%v", result["state"]),
			Status:      fmt.Sprintf("%v", result["status"]),
			DataType:    fmt.Sprintf("%v", result["data_type"]),
			Value:       fmt.Sprintf("%v", result["value"]),
			Threshold:   fmt.Sprintf("%v", result["threshold"]),
		}
		properties = append(properties, property)
	}

	return resterrors.SendOK(c, properties)

}

func (dc *DeviceController) getDeviceData(deviceId string) (models.Device, error) {
	device := models.Device{}
	query := `
				SELECT
					d.*,
					p.id AS property_id,
					p.reference_id AS referenceId,
					p.name AS name,
					p.unit AS unit,
					p.state AS state,
					p.status AS status,
					p.data_type AS dataType,
					p.value AS value,
					p.threshold AS threshold
				FROM
					devices d
				LEFT JOIN
					properties p ON d.id = p.device_id
				WHERE
					d.id = $1
			`

	result, queryErr := dc.repository.ExecuteQuery(query, deviceId)
	if queryErr != nil {
		log.Error("Error while executing query: %v", queryErr)
		return models.Device{}, error(queryErr)
	}
	log.Debug("Result: %v", result)
	if len(result) > 0 {
		device.ID = result[0]["id"].(string)
		device.ReferenceID = result[0]["reference_id"].(string)
		device.Type = result[0]["type"].(string)
		device.DeviceName = result[0]["device_name"].(string)
		device.CreatedAt = result[0]["created_at"].(string)
		device.State = result[0]["state"].(string)
		device.Location = result[0]["location"].(string)
		device.Status = result[0]["status"].(string)
		device.Customer = result[0]["customer"].(string)
		device.Site = result[0]["site"].(string)

		for _, propMap := range result {
			prop := models.Property{
				ID:          propMap["property_id"].(string),
				ReferenceID: propMap["reference_id"].(string),
				Name:        propMap["name"].(string),
				Unit:        propMap["unit"].(string),
				State:       propMap["state"].(string),
				Status:      propMap["status"].(string),
				DataType:    propMap["datatype"].(string),
				Value:       propMap["value"].(string),
				Threshold:   propMap["threshold"].(string),
			}
			device.Properties = append(device.Properties, prop)
		}
	}

	return device, nil
}

func (dc *DeviceController) GetAllDevices(c *fiber.Ctx) error {
	// Use the repository to get all devices

	deviceList := []models.Device{}
	deviceIds, err := dc.repository.GetAll("devices")
	if err != nil {
		return resterrors.SendInternalServerError(c)
	}

	log.Debug("Devices: %v", deviceIds)
	for _, device := range deviceIds {
		log.Debug("Device ID: %s", device["id"])
		deviceId := device["id"].(string)

		device, err := dc.getDeviceData(deviceId)
		if err != nil {
			return resterrors.SendInternalServerError(c)
		}

		// devices = append(devices, device)
		if !reflect.DeepEqual(device, models.Device{}) {
			deviceList = append(deviceList, device)
		}
	}

	// if len(deviceList) == 0 {
	// 	return resterrors.SendNotFoundError(c, "No devices found")
	// }

	return resterrors.SendOK(c, deviceList)
}

func (dc *DeviceController) GetDeviceByID(c *fiber.Ctx) error {
	// Get the device ID from the request parameters
	id := c.Params("id")

	// Use the repository to get the device by ID
	device, err := dc.getDeviceData(id)
	if err != nil {
		return resterrors.SendNotFoundError(c, "No device found")
	}
	if reflect.DeepEqual(device, models.Device{}) {
		return resterrors.SendNotFoundError(c, "No device found")
	}

	return resterrors.SendOK(c, device)
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
func (dc *DeviceController) CreateDevice(c *fiber.Ctx) error {
	log.Info("Creating a new device")
	var device models.Device
	err := c.BodyParser(&device)
	log.Debug("Err: %e", err)
	if err != nil {
		return resterrors.SendBadRequestError(c)
	}
	// Validate the device
	validate := validator.New()
	err = validate.Struct(device)
	if err != nil {
		return resterrors.SendBadRequestError(c)
	}
	log.Debug("Device: %v", device)
	deviceData, getErr := dc.repository.CheckExist("devices", "reference_id", device.ReferenceID)
	if getErr != nil {
		return resterrors.SendInternalServerError(c)
	}
	if deviceData {
		return resterrors.SendConflictError(c, "Device already exists")
	}
	log.Debug("Device Data: %v", deviceData)
	deviceMap := make(map[string]interface{})
	deviceMap["reference_id"] = device.ReferenceID
	deviceMap["type"] = device.Type
	deviceMap["device_name"] = device.DeviceName
	deviceMap["created_at"] = device.CreatedAt
	deviceMap["state"] = device.State
	deviceMap["location"] = device.Location
	deviceMap["status"] = device.Status
	deviceMap["customer"] = device.Customer
	deviceMap["site"] = device.Site
	// Use the repository to create the device
	log.Debug("Device Map: %v", deviceMap)

	devId, devErr := dc.repository.Create("devices", deviceMap)
	if devErr != nil {
		return resterrors.SendInternalServerError(c)
	}

	var propertiesMap []map[string]interface{}
	for _, prop := range device.Properties {
		propMap := make(map[string]interface{})
		propMap["id"] = prop.ID
		propMap["reference_id"] = prop.ReferenceID
		propMap["name"] = prop.Name
		propMap["unit"] = prop.Unit
		propMap["state"] = prop.State
		propMap["status"] = prop.Status
		propMap["data_type"] = prop.DataType
		propMap["value"] = prop.Value
		propMap["threshold"] = prop.Threshold
		propMap["device_id"] = devId
		propId, propErr := dc.repository.Create("properties", propMap)
		if propErr != nil {
			return resterrors.SendInternalServerError(c)
		}
		log.Debug("Property ID: %s", propId)
		propertiesMap = append(propertiesMap, propMap)
	}
	deviceMap["properties"] = propertiesMap
	log.Debug("Device Map: %v", deviceMap)
	log.Debug("Properties Map: %v", propertiesMap)

	// Create a new Device object
	device = models.Device{
		ID:          devId,
		ReferenceID: deviceMap["reference_id"].(string),
		Type:        deviceMap["type"].(string),
		DeviceName:  deviceMap["device_name"].(string),
		CreatedAt:   deviceMap["created_at"].(string),
		State:       deviceMap["state"].(string),
		Location:    deviceMap["location"].(string),
		Status:      deviceMap["status"].(string),
		Customer:    deviceMap["customer"].(string),
		Site:        deviceMap["site"].(string),
	}

	// Create a slice of Property objects
	properties := make([]models.Property, 0)
	for _, propMap := range propertiesMap {
		prop := models.Property{
			ID:          propMap["id"].(string),
			ReferenceID: propMap["reference_id"].(string),
			Name:        propMap["name"].(string),
			Unit:        propMap["unit"].(string),
			State:       propMap["state"].(string),
			Status:      propMap["status"].(string),
			DataType:    propMap["data_type"].(string),
			Value:       propMap["value"].(string),
			Threshold:   propMap["threshold"].(string),
		}
		properties = append(properties, prop)
	}
	device.Properties = properties
	// fmt.Println(propertiesMap)
	return resterrors.SendCreated(c, device)
}

func (dc *DeviceController) UpdateDevice(c *fiber.Ctx) error {
	// Get the device ID from the request parameters
	id := c.Params("id")

	// Create a new device instance from the request body
	var device map[string]interface{}
	err := c.BodyParser(&device)
	if err != nil {
		return resterrors.SendBadRequestError(c)
	}

	// Use the repository to update the device
	affectedRows, err := dc.repository.Update("devices", id, device)
	if err != nil {
		return resterrors.SendInternalServerError(c)
	}

	if affectedRows == 0 {
		return resterrors.SendNotFoundError(c, "No device found")
	}

	return resterrors.SendOK(c, affectedRows)
}

func (dc *DeviceController) DeleteDevice(c *fiber.Ctx) error {
	// Get the device ID from the request parameters
	id := c.Params("id")

	// Use the repository to delete the device
	deviceData, getErr := dc.repository.CheckExist("devices", "id", id)
	if getErr != nil {
		return resterrors.SendInternalServerError(c)
	}
	if !deviceData {
		return resterrors.SendNotFoundError(c, "Device not found")
	}
	propErr := dc.repository.Delete("properties", "device_id", id)
	log.Debug("Delete Error %v", propErr)
	if propErr != nil {
		return resterrors.SendInternalServerError(c)
	}

	err := dc.repository.Delete("devices", "id", id)
	log.Debug("Delete Error %v", err)
	if err != nil {
		return resterrors.SendInternalServerError(c)
	}
	return resterrors.SendNoContent(c)
}
