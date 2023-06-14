package app

import (
	"context"
	"fmt"
	"io"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/RamazanZholdas/KeyboardistSV2/pkg/database"
	"github.com/RamazanZholdas/KeyboardistSV2/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type App struct {
	Fiber *fiber.App
}

var (
	mongoInstance = &database.MongoDB{}
	fiberLogFile  *os.File
)

func Intitialize(mongoURI, dbName string) (*App, error) {
	utils.CreateLogFiles()

	err := mongoInstance.Connect(mongoURI, dbName)
	if err != nil {
		utils.LogError("Error connecting to MongoDB: ", err)
		return nil, fmt.Errorf("error connecting to MongoDB: %v", err)
	}

	var allCollections = []string{
		os.Getenv("COLLECTION_USERS"),
		os.Getenv("COLLECTION_PRODUCTS"),
		os.Getenv("COLLECTION_NEWS"),
		os.Getenv("COLLECTION_WIKI"),
		os.Getenv("COLLECTION_SERVICE_MASTERS"),
		os.Getenv("COLLECTION_PURCHASE_HISTORY"),
		os.Getenv("COLLECTION_POSTS"),
		os.Getenv("COLLECTION_COMMENTS"),
	}

	var wg sync.WaitGroup
	wg.Add(len(allCollections))
	errors := make(chan error, len(allCollections))

	for _, collectionName := range allCollections {
		go func(name string) {
			defer wg.Done()
			err := mongoInstance.CreateCollection(name)
			if err != nil {
				if !strings.Contains(err.Error(), "already exists") {
					errors <- fmt.Errorf("error creating collection %s: %v", name, err)
				} else {
					utils.LogWarning(fmt.Sprintf("Collection %s already exists", name))
				}
			}
		}(collectionName)
	}

	go func() {
		wg.Wait()
		close(errors)
	}()

	for err := range errors {
		utils.LogError(err)
		return nil, err
	}

	coll := mongoInstance.Db.Collection(os.Getenv("COLLECTION_PRODUCTS"))
	model := mongo.IndexModel{Keys: bson.D{{Key: "name", Value: "text"}}}
	name, err := coll.Indexes().CreateOne(context.TODO(), model)
	if err != nil {
		utils.LogError("Error creating index: ", err)
		return nil, fmt.Errorf("error creating index: %v", err)
	}

	utils.LogInfo(fmt.Sprintf("Created index %s", name))

	fiberLogFile, err = os.OpenFile("./logs/FiberLogs.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		utils.LogError("Error opening fiberlog file: ", err)
		return nil, fmt.Errorf("error opening fiberlog file: %v", err)
	}

	fiber := fiber.New()

	fiber.Use(cors.New(cors.Config{
		AllowCredentials: true,
	}))

	mw := io.MultiWriter(os.Stdout, fiberLogFile)

	fiber.Use(logger.New(logger.Config{
		Format:     "[${ip}]:${port} ${status} - ${method} ${path}\n",
		Output:     mw,
		TimeFormat: time.RFC3339Nano,
		TimeZone:   "Local",
	}))

	return &App{
		Fiber: fiber,
	}, nil
}

func (a *App) Close() {
	mongoInstance.Disconnect()
	a.Fiber.Shutdown()
	fiberLogFile.Close()
	utils.LogInfo("Closing app")
	utils.CloseLogFiles()
}

func GetMongoInstance() *database.MongoDB {
	return mongoInstance
}
