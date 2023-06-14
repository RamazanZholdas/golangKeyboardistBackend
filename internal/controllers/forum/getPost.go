package forum

import (
	"github.com/RamazanZholdas/KeyboardistSV2/internal/app"
	"github.com/RamazanZholdas/KeyboardistSV2/internal/models"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"os"
	"strconv"
	"time"
)

type responseBody struct {
	ID           primitive.ObjectID                            `json:"id"`
	Order        int32                                         `json:"order"`
	Title        string                                        `json:"title"`
	Content      string                                        `json:"content"`
	Author       string                                        `json:"author"`
	Date         string                                        `json:"date"`
	LastModified string                                        `json:"last_modified"`
	DayAgo       int                                           `json:"day_ago"`
	Bounder      []map[primitive.ObjectID][]primitive.ObjectID `json:"bounder"`
	MainComment  []models.Comment                              `json:"main_comment"`
	ReplyComment []models.Comment                              `json:"reply_comment"`
	//Comments     []map[models.Comment][]models.Comment         `json:"comments"`
}

/*
	[
		{
			Maincomment1: [replycomment1, replycomment1, replycomment1]
		},
		{
			Maincomment2: [replycomment2, replycomment2, replycomment2]
		},
	]

*/

func GetPost(c *fiber.Ctx) error {
	Order := c.Params("order")
	number, err := strconv.Atoi(Order)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"message": "The order must be a valid number.",
		})
	}

	var post models.Post
	filter := bson.M{"order": number}
	err = app.GetMongoInstance().FindOne(os.Getenv("COLLECTION_POSTS"), filter, &post)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{
			"message": "Product not found.",
		})
	}

	var body responseBody

	for _, comment := range post.Comments {
		for mc, rc := range comment {
			var mainComment models.Comment
			var replyComments []models.Comment

			err = app.GetMongoInstance().FindOne(os.Getenv("COLLECTION_COMMENTS"), bson.M{"_id": mc}, &mainComment)
			if err != nil {
				return c.Status(404).JSON(fiber.Map{
					"message": "Comment not found.",
				})
			}

			//if len(rc) == 0 {

			for _, reply := range rc {
				var replyComment models.Comment

				err = app.GetMongoInstance().FindOne(os.Getenv("COLLECTION_COMMENTS"), bson.M{"_id": reply}, &replyComment)
				if err != nil {
					return c.Status(404).JSON(fiber.Map{
						"message": "Comment not found.",
					})
				}

				replyComments = append(replyComments, replyComment)
				body.ReplyComment = append(body.ReplyComment, replyComment)
			}
			body.MainComment = append(body.MainComment, mainComment)
			temp := map[models.Comment][]models.Comment{}
			temp[mainComment] = replyComments
			mainComment = models.Comment{}
			replyComments = []models.Comment{}
		}
	}
	postDate, _ := time.Parse("2006-01-02 15:04:05", post.Date)
	daysAgo := int(time.Since(postDate).Hours() / 24)

	body.ID = post.ID
	body.Order = post.Order
	body.Title = post.Title
	body.Content = post.Content
	body.Author = post.Author
	body.Date = post.Date
	body.LastModified = post.LastModified
	body.DayAgo = daysAgo
	//body.Comments = tempBodyComment
	body.Bounder = post.Comments

	//fmt.Println(body.Comments)

	return c.JSON(body)
}
