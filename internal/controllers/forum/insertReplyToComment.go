package forum

import (
	"encoding/json"
	"github.com/RamazanZholdas/KeyboardistSV2/internal/app"
	"github.com/RamazanZholdas/KeyboardistSV2/internal/jwt"
	"github.com/RamazanZholdas/KeyboardistSV2/internal/models"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"os"
	"strconv"
	"time"
)

/*
	json body:
	{
		"order": "order",
		"content": "content",
		"comment_id": "comment_id"
	}

	TODO: return post with updated comments
*/

func InsertReplyToComment(c *fiber.Ctx) error {
	cookie := c.Cookies("jwt")

	claims, err := jwt.ExtractTokenClaimsFromCookie(cookie)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"message": "Unauthorized"})
	}

	var requestBody struct {
		Order     string `json:"order"`
		Content   string `json:"content"`
		CommentId string `json:"comment_id"`
	}

	if err := json.Unmarshal(c.Body(), &requestBody); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Invalid request body"})
	}

	var post models.Post
	number, err := strconv.Atoi(requestBody.Order)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"message": "The order must be a valid number.",
		})
	}
	post_order := number
	err = app.GetMongoInstance().FindOne(os.Getenv("COLLECTION_POSTS"), bson.M{"order": int32(post_order)}, &post)
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

	filter := bson.M{"date": comment.Date}
	err = app.GetMongoInstance().FindOne(os.Getenv("COLLECTION_COMMENTS"), filter, &comment)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{
			"message": "Comment not found.",
		})
	}

	comment_oid, err := primitive.ObjectIDFromHex(requestBody.CommentId)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"message": "Invalid post ID.",
		})
	}

	for _, coms := range post.Comments {
		if findVariableInMap(comment_oid, coms) {
			coms[comment_oid] = append(coms[comment_oid], comment.ID)
		}
	}

	post.LastModified = time.Now().Format("2006-01-02 15:04:05")
	filter = bson.M{"_id": post.ID}

	err = app.GetMongoInstance().UpdateOne(os.Getenv("COLLECTION_POSTS"), filter, bson.M{"$set": bson.M{"comments": post.Comments}})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Internal server error"})
	}

	err = app.GetMongoInstance().UpdateOne(os.Getenv("COLLECTION_POSTS"), filter, bson.M{"$set": bson.M{"last_modified": post.LastModified}})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Internal server error"})
	}

	return c.Status(200).JSON(fiber.Map{
		"message": "Comment added successfully.",
	})
}

func findVariableInMap(variable primitive.ObjectID, m map[primitive.ObjectID][]primitive.ObjectID) bool {
	_, ok := m[variable]
	return ok
}
