package news

import (
	"os"

	"github.com/RamazanZholdas/KeyboardistSV2/internal/app"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
)

func GetAllNews(c *fiber.Ctx) error {
	var news []bson.M
	err := app.GetMongoInstance().FindMany(os.Getenv("COLLECTION_NEWS"), bson.M{}, &news)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"message": "Internal server error",
		})
	}

	if len(news) == 0 {
		return c.Status(404).JSON(fiber.Map{
			"message": "No products found",
		})
	}

	return c.JSON(news)
}
