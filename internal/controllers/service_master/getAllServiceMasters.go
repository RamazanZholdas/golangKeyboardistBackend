package servicemaster

import (
	"os"

	"github.com/RamazanZholdas/KeyboardistSV2/internal/app"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
)

func GetAllServiceMasters(c *fiber.Ctx) error {
	var service_master []bson.M
	err := app.GetMongoInstance().FindMany(os.Getenv("COLLECTION_SERVICE_MASTERS"), bson.M{}, &service_master)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"message": "Internal server error",
		})
	}

	if len(service_master) == 0 {
		return c.Status(404).JSON(fiber.Map{
			"message": "No service master found",
		})
	}

	return c.JSON(service_master)
}
