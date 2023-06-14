package forum

import (
	"fmt"
	"os"
	"sort"
	"time"

	"github.com/RamazanZholdas/KeyboardistSV2/internal/app"
	"github.com/RamazanZholdas/KeyboardistSV2/internal/models"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func GetAllPosts(c *fiber.Ctx) error {
	var posts []primitive.M
	err := app.GetMongoInstance().FindMany(os.Getenv("COLLECTION_POSTS"), bson.M{}, &posts)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"message": "Internal server error",
		})
	}

	if len(posts) == 0 {
		empty := []string{}
		return c.JSON(empty)
	}

	var postModels []models.Post
	for _, post := range posts {
		postDate, _ := time.Parse("2006-01-02 15:04:05", post["date"].(string))
		daysAgo := int(time.Since(postDate).Hours() / 24)
		var daysAgoStr string
		if daysAgo == 0 {
			daysAgoStr = "today"
		} else {
			daysAgoStr = fmt.Sprint(daysAgo) + " days ago"
		}

		var content_for_response string
		if len(post["date"].(string)) > 50 {
			content_for_response = post["content"].(string)[:50] + "..."
		} else {
			content_for_response = post["content"].(string)
		}

		postModel := models.Post{
			ID:      post["_id"].(primitive.ObjectID),
			Title:   post["title"].(string),
			Content: content_for_response,
			Author:  post["author"].(string),
			Date:    post["date"].(string),
			Order:   post["order"].(int32),
			DayAgo:  daysAgoStr,
		}
		postModels = append(postModels, postModel)
	}

	sort.Slice(postModels, func(i, j int) bool {
		ti, err := time.Parse("2006-01-02 15:04:05", postModels[i].Date)
		if err != nil {
			return false
		}
		tj, err := time.Parse("2006-01-02 15:04:05", postModels[j].Date)
		if err != nil {
			return true
		}
		return ti.After(tj)
	})

	return c.JSON(postModels)
}
