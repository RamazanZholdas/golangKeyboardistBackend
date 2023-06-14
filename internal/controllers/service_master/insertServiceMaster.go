package servicemaster

import (
	"os"

	"github.com/RamazanZholdas/KeyboardistSV2/internal/app"
	"github.com/RamazanZholdas/KeyboardistSV2/internal/models"
	"github.com/gofiber/fiber/v2"
)

func InsertServiceMaster(c *fiber.Ctx) error {
	var service_master models.ServiceMaster

	if err := c.BodyParser(&service_master); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"message": "Bad Request",
		})
	}

	count, err := app.GetMongoInstance().CountDocuments(os.Getenv("COLLECTION_SERVICE_MASTERS"))
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"message": "Internal server error",
		})
	}
	count++

	service_master.Order = int32(count)

	insertErr := app.GetMongoInstance().InsertOne(os.Getenv("COLLECTION_SERVICE_MASTERS"), service_master)
	if insertErr != nil {
		return c.Status(500).JSON(fiber.Map{
			"message": "Internal server error",
		})
	}

	return c.JSON(service_master)
}
