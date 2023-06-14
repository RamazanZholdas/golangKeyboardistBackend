package forum

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/RamazanZholdas/KeyboardistSV2/internal/app"
	"github.com/RamazanZholdas/KeyboardistSV2/internal/jwt"
	"github.com/RamazanZholdas/KeyboardistSV2/internal/models"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

/*
	json body:
	{
		"order": "order",
		"content": "content",
	}

	TODO:return post with updated comments
*/

func InsertComment(c *fiber.Ctx) error {
	cookie := c.Cookies("jwt")

	claims, err := jwt.ExtractTokenClaimsFromCookie(cookie)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"message": "Unauthorized"})
	}

	var requestBody struct {
		Order   string `json:"order"`
		Content string `json:"content"`
	}

	if err := json.Unmarshal(c.Body(), &requestBody); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Invalid request body"})
	}

	var post models.Post
	post_order := requestBody.Order
	number, err := strconv.Atoi(post_order)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"message": "The order must be a valid number.",
		})
	}
	filter := bson.M{"order": number}

	err = app.GetMongoInstance().FindOne(os.Getenv("COLLECTION_POSTS"), filter, &post)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"message": "Post not found"})
	}

	var comment models.Comment
	comment.Author = claims.Issuer
	comment.Content = requestBody.Content
	comment.Date = time.Now().Format("2006-01-02 15:04:05")

	err = app.GetMongoInstance().InsertOne(os.Getenv("COLLECTION_COMMENTS"), comment)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Internal server error"})
	}

	filter = bson.M{"date": comment.Date}
	err = app.GetMongoInstance().FindOne(os.Getenv("COLLECTION_COMMENTS"), filter, &comment)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{
			"message": "Comment not found.",
		})
	}

	var postComment []map[primitive.ObjectID][]primitive.ObjectID
	emptySlice := make([]primitive.ObjectID, 0)
	tempComment := make(map[primitive.ObjectID][]primitive.ObjectID)
	tempComment[comment.ID] = emptySlice
	postComment = append(postComment, tempComment)

	post.Comments = append(post.Comments, postComment...)
	post.LastModified = time.Now().Format("2006-01-02 15:04:05")
	filter = bson.M{"_id": post.ID}

	fmt.Println(post.Comments)

	err = app.GetMongoInstance().UpdateOne(os.Getenv("COLLECTION_POSTS"), filter, bson.M{"$set": bson.M{"comments": post.Comments}})
	if err != nil {
		fmt.Println("I cant put it in the post")
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Internal server error"})
	}

	err = app.GetMongoInstance().UpdateOne(os.Getenv("COLLECTION_POSTS"), filter, bson.M{"$set": bson.M{"last_modified": post.LastModified}})
	if err != nil {
		fmt.Println("Cant save last modified")
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Internal server error"})
	}

	return c.Status(200).JSON(fiber.Map{
		"message": "Comment added successfully.",
	})
}
