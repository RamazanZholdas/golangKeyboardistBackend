package product

import (
	"os"

	"github.com/RamazanZholdas/KeyboardistSV2/internal/app"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func FindProduct(c *fiber.Ctx) error {
	var productName = c.Params("productName")
	if productName == "" {
		return c.Status(400).JSON(fiber.Map{
			"message": "Bad request",
		})
	}

	var products []primitive.M
	filter := bson.D{{Key: "$text", Value: bson.D{{Key: "$search", Value: "\"" + productName + "\""}}}}

	err := app.GetMongoInstance().FindMany(os.Getenv("COLLECTION_PRODUCTS"), filter, &products)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"message": "Internal server error",
		})
	}

	if len(products) == 0 {
		empty := []string{}
		return c.JSON(empty)
	}

	return c.JSON(products)
}
