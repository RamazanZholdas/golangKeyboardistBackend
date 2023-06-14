package product

import (
	"os"

	"github.com/RamazanZholdas/KeyboardistSV2/internal/app"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func GetAllProducts(c *fiber.Ctx) error {
	var products []primitive.M
	err := app.GetMongoInstance().FindMany(os.Getenv("COLLECTION_PRODUCTS"), bson.M{}, &products)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"message": "Internal server error",
		})
	}

	if len(products) == 0 {
		return c.Status(404).JSON(fiber.Map{
			"message": "No products found",
		})
	}

	return c.JSON(products)
}
