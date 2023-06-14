package cart

import (
	"encoding/json"
	"os"
	"time"

	"github.com/RamazanZholdas/KeyboardistSV2/internal/app"
	"github.com/RamazanZholdas/KeyboardistSV2/internal/jwt"
	"github.com/RamazanZholdas/KeyboardistSV2/internal/models"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
)

func PayNow(c *fiber.Ctx) error {
	var requestBody struct {
		TotalPrice string `json:"totalPrice"`
	}

	if err := json.Unmarshal(c.Body(), &requestBody); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Invalid request body"})
	}

	totalPrice := requestBody.TotalPrice

	cookie := c.Cookies("jwt")

	claims, err := jwt.ExtractTokenClaimsFromCookie(cookie)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"message": "Unauthorized"})
	}

	var user models.User
	err = app.GetMongoInstance().FindOne("users", bson.M{"email": claims.Issuer}, &user)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"message": "User not found"})
	}

	purchaseHistory := models.PurchaseHistory{
		UserID:       user.ID,
		Products:     user.Cart,
		TotalPrice:   totalPrice,
		PurchaseDate: time.Now().Format("2006-01-02 15:04:05"),
		PaymentType:  "PayPal",
	}

	cart := []map[string]models.Product{}
	user.Cart = cart

	err = app.GetMongoInstance().UpdateOne("users", bson.M{"email": claims.Issuer}, bson.M{"$set": bson.M{"cart": user.Cart}})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Internal server error"})
	}

	insertErr := app.GetMongoInstance().InsertOne(os.Getenv("COLLECTION_PURCHASE_HISTORY"), purchaseHistory)
	if insertErr != nil {
		return c.Status(500).JSON(fiber.Map{
			"message": "Internal server error",
		})
	}

	return c.JSON(fiber.Map{
		"message": "success",
	})
}
