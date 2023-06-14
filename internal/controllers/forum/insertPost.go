package forum

import (
	"os"
	"time"

	"github.com/RamazanZholdas/KeyboardistSV2/internal/app"
	"github.com/RamazanZholdas/KeyboardistSV2/internal/jwt"
	"github.com/RamazanZholdas/KeyboardistSV2/internal/models"
	"github.com/gofiber/fiber/v2"
)

/*
json body:

	{
		"content": "content",
		"title": "title"
	}
*/
func InsertPost(c *fiber.Ctx) error {
	cookie := c.Cookies("jwt")

	claims, err := jwt.ExtractTokenClaimsFromCookie(cookie)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"message": "Unauthorized"})
	}

	var post models.Post

	if err := c.BodyParser(&post); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"message": "Bad Request",
		})
	}

	post.Author = claims.Issuer
	post.Date = time.Now().Format("2006-01-02 15:04:05")
	post.LastModified = time.Now().Format("2006-01-02 15:04:05")

	count, err := app.GetMongoInstance().CountDocuments(os.Getenv("COLLECTION_POSTS"))
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"message": "Internal server error",
		})
	}
	count++

	post.Order = int32(count)

	insertErr := app.GetMongoInstance().InsertOne(os.Getenv("COLLECTION_POSTS"), post)
	if insertErr != nil {
		return c.Status(500).JSON(fiber.Map{
			"message": "Internal server error",
		})
	}

	return c.JSON(post)
}
