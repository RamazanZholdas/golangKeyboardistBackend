package cart

import (
	"encoding/json"
	"strconv"

	"github.com/RamazanZholdas/KeyboardistSV2/internal/app"
	"github.com/RamazanZholdas/KeyboardistSV2/internal/jwt"
	"github.com/RamazanZholdas/KeyboardistSV2/internal/models"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
)

func ChangeQuantity(c *fiber.Ctx) error {
	var requestBody struct {
		Quantity string `json:"quantity"`
		OptionId string `json:"optionId"`
	}

	if err := json.Unmarshal(c.Body(), &requestBody); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Invalid request body"})
	}

	quantity := requestBody.Quantity
	optionId := requestBody.OptionId

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

	order := c.Params("order")
	number, err := strconv.Atoi(order)

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "The order must be a valid number.",
		})
	}

	var product models.Product
	err = app.GetMongoInstance().FindOne("products", bson.M{"order": number}, &product)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"message": "Product not found"})
	}

	for index, item := range user.Cart {
		if item["product"].Order == int32(number) {
			if item["product"].Options[0]["optionId"] != optionId {
				continue
			}

			if item["product"].Options[0]["inStock"] < quantity {
				return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
					"message": "The quantity must be less than or equal to the number of items in stock.",
				})
			}

			oldQuantity, err := strconv.Atoi(item["product"].Options[0]["quantity"])
			if err != nil {
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Internal server error"})
			}
			newQuantity, err := strconv.Atoi(quantity)
			if err != nil {
				return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Request body is invalid"})
			}

			resultQuantity := oldQuantity + newQuantity

			if resultQuantity == 0 {
				cart := []map[string]models.Product{}
				for _, item := range user.Cart {
					if item["product"].Order == int32(number) {
						cart = append(user.Cart[:index], user.Cart[index+1:]...)
					}
				}

				err = app.GetMongoInstance().UpdateOne("users", bson.M{"email": claims.Issuer}, bson.M{"$set": bson.M{"cart": cart}})
				if err != nil {
					return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Internal server error"})
				}

				return c.Status(fiber.StatusOK).JSON(fiber.Map{
					"message": "The product has been removed from the cart successfully.",
				})
			}

			item["product"].Options[0]["quantity"] = strconv.Itoa(resultQuantity)

			user.Cart[index]["product"] = item["product"]

			err = app.GetMongoInstance().UpdateOne("users", bson.M{"email": claims.Issuer}, bson.M{"$set": bson.M{"cart": user.Cart}})
			if err != nil {
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Internal server error"})
			}
			break
		}
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "The quantity of the product has been changed successfully.",
	})
}
