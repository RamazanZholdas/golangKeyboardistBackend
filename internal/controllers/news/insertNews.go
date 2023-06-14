package news

import (
	"os"
	"time"

	"github.com/RamazanZholdas/KeyboardistSV2/internal/app"
	"github.com/RamazanZholdas/KeyboardistSV2/internal/models"
	"github.com/gofiber/fiber/v2"
)

func InsertNews(c *fiber.Ctx) error {
	var news models.News

	if err := c.BodyParser(&news); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"message": "Bad Request",
		})
	}

	count, err := app.GetMongoInstance().CountDocuments(os.Getenv("COLLECTION_NEWS"))
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"message": "Internal server error",
		})
	}
	count++

	news.Order = int32(count)
	news.Date = string(time.Now().Format("02/01/2006"))

	insertErr := app.GetMongoInstance().InsertOne(os.Getenv("COLLECTION_NEWS"), news)
	if insertErr != nil {
		return c.Status(500).JSON(fiber.Map{
			"message": "Internal server error",
		})
	}

	return c.JSON(news)
}
