package news

import (
	"os"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"

	"github.com/RamazanZholdas/KeyboardistSV2/internal/app"
)

func GetNews(c *fiber.Ctx) error {
	Order := c.Params("order")
	number, err := strconv.Atoi(Order)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"message": "The order must be a valid number.",
		})
	}

	var news bson.M
	filter := bson.M{"order": number}

	err = app.GetMongoInstance().FindOne(os.Getenv("COLLECTION_NEWS"), filter, &news)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{
			"message": "Product not found.",
		})
	}

	return c.JSON(news)
}
