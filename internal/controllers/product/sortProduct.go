package product

import (
	"os"

	"github.com/RamazanZholdas/KeyboardistSV2/internal/app"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
)

func SortProduct(c *fiber.Ctx) error {
	var sortType = c.Params("sortType")
	if sortType != "price_asc" && sortType != "price_desc" {
		return c.Status(400).JSON(fiber.Map{
			"message": "Bad request",
		})
	}

	sortOptions := newSortOptions(sortType)

	products, err := app.GetMongoInstance().SortedDocuments(sortOptions, os.Getenv("COLLECTION_PRODUCTS"))
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"message": "Internal server error",
		})
	}

	return c.JSON(products)
}

func newSortOptions(sortType string) bson.M {
	var sortOptions bson.M

	sortOptions = bson.M{"options.0.price": 1}
	if sortType == "price_desc" {
		sortOptions = bson.M{"options.0.price": -1}
	}

	return sortOptions
}
