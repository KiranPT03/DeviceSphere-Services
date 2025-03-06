package devices

import (

	// Import the PostgreSQL repository package
	"encoding/json"
	"fmt"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	config "github.com/kiranpt03/factorysphere/devicesphere/services/deviceservices/pkg/config"
	postgres "github.com/kiranpt03/factorysphere/devicesphere/services/deviceservices/pkg/databases/postgres"
	models "github.com/kiranpt03/factorysphere/devicesphere/services/deviceservices/pkg/models"
	"github.com/kiranpt03/factorysphere/devicesphere/services/deviceservices/pkg/utils/resterrors"
)

type DeviceDiscoveryController struct {
	repository *postgres.PostgreSQLRepository
}

func NewDeviceDiscoveryController(config *config.Config) *DeviceDiscoveryController {
	// Create a new PostgreSQL repository instance
	repository, err := postgres.NewPostgreSQLRepository(config)
	if err != nil {
		panic(err)
	}

	return &DeviceDiscoveryController{repository: repository}
}

func (ddc *DeviceDiscoveryController) GetKepwareDeviceData(c *fiber.Ctx) error {
	filePath := "./pkg/files/kepwareDeviceData.json" // Replace with your file path
	file, err := os.ReadFile(filePath)
	if err != nil {
		fmt.Println("Error reading file:", err)
		return c.Status(fiber.StatusInternalServerError).SendString("Error reading file")
	}

	var channels []models.Channel
	err = json.Unmarshal(file, &channels)
	if err != nil {
		fmt.Println("Error unmarshaling JSON:", err)
		return c.Status(fiber.StatusInternalServerError).SendString("Error unmarshaling JSON")
	}

	var devices []models.Device

	for _, channel := range channels {
		for _, deviceData := range channel.Devices {
			device := models.Device{
				ID:          uuid.New().String(),
				ReferenceID: uuid.New().String(),
				Type:        channel.ChannelName,
				DeviceName:  deviceData.DeviceName,
				State:       "Active",
				Status:      "Operational",
				Properties:  []models.Property{},
			}

			for _, tag := range deviceData.Tags {
				property := models.Property{
					ID:          uuid.New().String(),
					ReferenceID: tag.TagId,
					Name:        tag.TagName,
					State:       "Active",
					Status:      "Operational",
					Threshold:   "0",
				}
				device.Properties = append(device.Properties, property)
			}
			devices = append(devices, device)
		}
	}

	return resterrors.SendOK(c, devices)
}
