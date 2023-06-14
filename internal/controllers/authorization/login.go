package authorization

import (
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"

	"github.com/RamazanZholdas/KeyboardistSV2/internal/app"
	"github.com/RamazanZholdas/KeyboardistSV2/internal/jwt"
)

func Login(c *fiber.Ctx) error {
	var data map[string]string

	if err := c.BodyParser(&data); err != nil {
		return err
	}

	filter := bson.M{"email": data["email"]}

	var user bson.M

	err := app.GetMongoInstance().FindOne(os.Getenv("COLLECTION_USERS"), filter, &user)
	if err != nil {
		return c.Status(401).JSON(fiber.Map{
			"message": "User not found",
		})
	}

	password := user["password"].(primitive.Binary).Data

	if err := bcrypt.CompareHashAndPassword(password, []byte(data["password"])); err != nil {
		return c.Status(401).JSON(fiber.Map{
			"message": "Incorrect password",
		})
	}

	email := user["email"].(string)

	token, err := jwt.NewJwtTokenWithClaims(email)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"message": "Internal server error",
		})
	}

	cookie := fiber.Cookie{
		Name:     "jwt",
		Value:    token,
		Expires:  time.Now().Add(time.Hour * 24),
		HTTPOnly: true,
	}

	c.Cookie(&cookie)

	return c.JSON(fiber.Map{
		"message": "success",
	})
}
