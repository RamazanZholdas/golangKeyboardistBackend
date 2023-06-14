package wiki

import (
	"os"

	"github.com/RamazanZholdas/KeyboardistSV2/internal/app"
	"github.com/RamazanZholdas/KeyboardistSV2/internal/models"
	"github.com/gofiber/fiber/v2"
)

func InsertWikiPage(c *fiber.Ctx) error {
	var wiki models.Wiki

	if err := c.BodyParser(&wiki); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"message": "Bad Request",
		})
	}

	count, err := app.GetMongoInstance().CountDocuments(os.Getenv("COLLECTION_WIKI"))
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"message": "Internal server error",
		})
	}
	count++

	wiki.Order = int32(count)

	insertErr := app.GetMongoInstance().InsertOne(os.Getenv("COLLECTION_WIKI"), wiki)
	if insertErr != nil {
		return c.Status(500).JSON(fiber.Map{
			"message": "Internal server error",
		})
	}

	return c.JSON(wiki)
}
