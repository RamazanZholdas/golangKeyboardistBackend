package cart

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/RamazanZholdas/KeyboardistSV2/internal/app"
	"github.com/RamazanZholdas/KeyboardistSV2/internal/jwt"
	"github.com/RamazanZholdas/KeyboardistSV2/internal/models"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
)

func InsertToCart(c *fiber.Ctx) error {
	var optionId string

	var requestBody struct {
		OptionId string `json:"optionId"`
	}

	if err := json.Unmarshal(c.Body(), &requestBody); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Invalid request body"})
	}
	optionId = requestBody.OptionId

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

	var userOptionList []string

	for _, cartItem := range user.Cart {
		userOptionList = append(userOptionList, cartItem["product"].Options[0]["optionId"])
	}

	for _, cartItem := range user.Cart {
		stringOrderCart := strconv.Itoa(int(cartItem["product"].Order))
		if order == stringOrderCart && cartItem["product"].Options[0]["optionId"] == optionId {
			fmt.Println("same product exists")
			fmt.Println(cartItem["product"].Options[0]["optionId"])
			fmt.Println(optionId)
			fmt.Println("same product, same option")
			quantity, err := strconv.Atoi(cartItem["product"].Options[0]["quantity"])
			if err != nil {
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Failed to update cart"})
			}

			inStock, err := strconv.Atoi(cartItem["product"].Options[0]["inStock"])
			if err != nil {
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Failed to update cart"})
			}

			if quantity == inStock {
				return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Product is out of stock"})
			}

			quantity++
			cartItem["product"].Options[0]["quantity"] = strconv.Itoa(quantity)

			update := bson.M{"$set": bson.M{"cart": user.Cart}}
			err = app.GetMongoInstance().UpdateOne("users", bson.M{"_id": user.ID}, update)
			if err != nil {
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Failed to update cart"})
			}
			return c.JSON(fiber.Map{"message": "Product added to cart"})

		} else if order == stringOrderCart && cartItem["product"].Options[0]["optionId"] != optionId && !contains(userOptionList, optionId) {
			fmt.Println("same product, different option")
			var productWithDifferentOpt models.Product
			err = app.GetMongoInstance().FindOne("products", bson.M{"order": number}, &productWithDifferentOpt)
			if err != nil {
				return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"message": "Product not found"})
			}

			filteredOptions := make([]map[string]string, 0)
			for _, option := range productWithDifferentOpt.Options {
				if option["optionId"] == optionId {
					filteredOptions = append(filteredOptions, option)
					break
				}
			}

			productWithDifferentOpt.Options = filteredOptions
			productWithDifferentOpt.Options[0]["quantity"] = "1"

			user.Cart = append(user.Cart, map[string]models.Product{"product": productWithDifferentOpt})

			update := bson.M{"$set": bson.M{"cart": user.Cart}}
			err = app.GetMongoInstance().UpdateOne("users", bson.M{"_id": user.ID}, update)
			if err != nil {
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Failed to update cart"})
			}

			return c.JSON(fiber.Map{"message": "Product added to cart"})
		}
	}

	var product models.Product
	err = app.GetMongoInstance().FindOne("products", bson.M{"order": number}, &product)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"message": "Product not found"})
	}

	filteredOptions := make([]map[string]string, 0)

	for _, option := range product.Options {
		if option["optionId"] == optionId {
			filteredOptions = append(filteredOptions, option)
			break
		}
	}

	product.Options = filteredOptions

	// setting up quantity
	product.Options[0]["quantity"] = "1"

	user.Cart = append(user.Cart, map[string]models.Product{"product": product})

	update := bson.M{"$set": bson.M{"cart": user.Cart}}
	err = app.GetMongoInstance().UpdateOne("users", bson.M{"_id": user.ID}, update)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Failed to update cart"})
	}

	return c.JSON(fiber.Map{"message": "Product added to cart"})
}

func contains(s []string, str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}

	return false
}
