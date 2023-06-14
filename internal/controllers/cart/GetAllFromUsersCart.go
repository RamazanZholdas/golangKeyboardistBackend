package cart

import (
	"github.com/RamazanZholdas/KeyboardistSV2/internal/app"
	"github.com/RamazanZholdas/KeyboardistSV2/internal/jwt"
	"github.com/RamazanZholdas/KeyboardistSV2/internal/models"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
)

func GetAllFromUsersCart(c *fiber.Ctx) error {
	cookie := c.Cookies("jwt")

	claims, err := jwt.ExtractTokenClaimsFromCookie(cookie)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid token"})
	}

	var user models.User
	err = app.GetMongoInstance().FindOne("users", bson.M{"email": claims.Issuer}, &user)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "User not found"})
	}

	if len(user.Cart) == 0 {
		var emptyProducts []models.Product
		return c.JSON(emptyProducts)
	}

	var products []models.Product
	for _, cartItem := range user.Cart {
		product := cartItem["product"]
		products = append(products, product)
	}

	return c.JSON(products)
}
