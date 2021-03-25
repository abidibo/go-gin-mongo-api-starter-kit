package database

import (
	"context"
	"fmt"
	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
	"sync"
	// "systems-management-api/core"
)

type DBDriver struct {
	D *mongo.Database
}

var db *DBDriver
var once sync.Once

// DB returns a singleton pointer to mongo.Database instance
func DB() *DBDriver {
	if db == nil {
		// check settings
		if !viper.IsSet("database.host") || !viper.IsSet("database.port") {
			zap.S().Fatal("Missing database host or port number")
		}
		DBHost := viper.GetString("database.host")
		DBPort := viper.GetInt("database.port")
		DBName := viper.GetString("database.dbName")

		once.Do(func() {
			ConnectionURI := fmt.Sprintf("%s:%d", DBHost, DBPort)
			// Set client options
			clientOptions := options.Client().ApplyURI(ConnectionURI)
			// Connect to MongoDB
			client, err := mongo.Connect(context.TODO(), clientOptions)
			if err != nil {
				zap.S().Fatal("Cannot connect to mongo db on ", ConnectionURI)
			} else {
				zap.S().Info("Succesfully connected to mongo @ ", ConnectionURI)
			}
			// Check the connection
			if client.Ping(context.TODO(), nil) != nil {
				zap.S().Fatal("Failed to ping mongo db")
			} else {
				zap.S().Debug("Ping to mongo was successful")
			}
			db = &DBDriver{client.Database(DBName)}
		})
	}

	return db
}

// Retrieves a single document given the collection and its ID
func (self *DBDriver) GetById(collectionName string, id string, item interface{}) error {
	var objId, err = primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	collection := self.D.Collection(collectionName)
	err = collection.FindOne(
		context.TODO(),
		bson.M{"_id": objId},
	).Decode(item)

	return err
}
