package database

import (
	"context"
	"errors"
	"log"

	"github.com/RamazanZholdas/KeyboardistSV2/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	ErrCreateDatabase = errors.New("failed to create database")
)

type MongoDB struct {
	Client *mongo.Client
	Db     *mongo.Database
}

func (m *MongoDB) Connect(uri, database string) error {
	utils.LogInfo("Connecting to MongoDB...")
	client, err := mongo.NewClient(options.Client().ApplyURI(uri))
	if err != nil {
		utils.LogError("Failed to create new MongoDB client: ", err)
		return errors.New("failed to create new MongoDB client")
	}

	err = client.Connect(context.Background())
	if err != nil {
		utils.LogError("Failed to connect to MongoDB: ", err)
		return errors.New("failed to connect to MongoDB")
	}

	m.Client = client
	m.Db = client.Database(database)

	databases, err := m.Client.ListDatabaseNames(context.Background(), bson.M{})
	if err != nil {
		utils.LogError("Failed to list database names: ", err)
		return ErrCreateDatabase
	}

	dbExists := false
	for _, db := range databases {
		if db == database {
			dbExists = true
			break
		}
	}

	if !dbExists {
		utils.LogWarning("Database does not exist, creating new database: ", database)
		err = m.Db.RunCommand(context.Background(), bson.D{{Key: "create", Value: database}}).Err()
		if err != nil {
			utils.LogError("Failed to create database: ", err)
			return ErrCreateDatabase
		}
	}

	utils.LogInfo("Successfully connected to MongoDB!")
	return nil
}

func (m *MongoDB) Disconnect() {
	err := m.Client.Disconnect(context.Background())
	if err != nil {
		utils.LogError("Failed to disconnect from MongoDB: ", err)
		log.Println("failed to disconnect from MongoDB")
	}
}
