package database

import (
	"context"
	"fmt"
	"time"

	"github.com/aiuzu42/AiuzuBotDiscord/config"
	"github.com/aiuzu42/AiuzuBotDiscord/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoDB struct {
	client     *mongo.Client
	collection *mongo.Collection
}

// InitDB start the mongoDB connection initializing the client and collection pointers.
func (m *MongoDB) InitDB(c config.DBConnection) *models.AppError {
	var err error
	options := options.Client().ApplyURI(fmt.Sprintf("mongodb+srv://%s:%s@%s/%s%s", c.User, c.Pass, c.Host, c.DB, c.Options))
	m.client, err = mongo.Connect(context.TODO(), options)
	if err != nil {
		return &models.AppError{Code: models.CantConnectToDatabase, Message: err.Error()}
	}
	err = m.client.Ping(context.TODO(), nil)
	if err != nil {
		return &models.AppError{Code: models.CantConnectToDatabase, Message: err.Error()}
	}
	m.collection = m.client.Database(c.DB).Collection(c.Collection)
	return nil
}

// GetUser searches and returns user information from the database that matches either the userId or username.
func (m *MongoDB) GetUser(userID string, username string) (models.User, *models.AppError) {
	query := bson.M{"$or": []bson.M{{"userID": userID}, {"fullName": username}}}
	var result models.User
	err := m.collection.FindOne(context.TODO(), query).Decode(&result)
	if err != nil && err == mongo.ErrNoDocuments {
		return models.User{}, &models.AppError{Code: models.UserNotFoundCode, Message: "User not found"}
	} else if err != nil {
		return models.User{}, &models.AppError{Code: models.DatabaseError, Message: err.Error()}
	}
	return result, nil
}

// AddUser adds the user information to the database.
func (m *MongoDB) AddUser(user models.User) *models.AppError {
	query := bson.M{
		"userID": user.UserID,
	}
	var found *mongo.SingleResult
	found = m.collection.FindOne(context.TODO(), query)
	if found.Err() == nil {
		return &models.AppError{Code: models.UserAlredyExists, Message: "The user alredy exists"}
	} else if found.Err() != nil && found.Err() != mongo.ErrNoDocuments {
		return &models.AppError{Code: models.DatabaseError, Message: found.Err().Error()}
	}
	_, err := m.collection.InsertOne(context.TODO(), user)
	if err != nil {
		return &models.AppError{Code: models.DatabaseError, Message: err.Error()}
	}
	return nil
}

// IncreaseMessageCount searches for the document with the provided userID.
// If found, it increases its server.messageCount value and updates server.lastMessage.
// If the user is not found an error is returned.
func (m *MongoDB) IncreaseMessageCount(userID string) *models.AppError {
	query := bson.M{
		"userID": userID,
	}
	var found *mongo.SingleResult
	found = m.collection.FindOne(context.TODO(), query)
	if found.Err() != nil && found.Err() == mongo.ErrNoDocuments {
		return &models.AppError{Code: models.UserNotFoundCode, Message: "User not found"}
	} else if found.Err() != nil {
		return &models.AppError{Code: models.DatabaseError, Message: found.Err().Error()}
	}
	lastMessage := time.Now().Format("01-02-2006")
	updateQuery := bson.D{
		{
			Key: "$inc", Value: bson.D{
				{Key: "server.messageCount", Value: 1},
			},
		},
		{
			Key: "$set", Value: bson.D{
				{Key: "server.lastMessage", Value: lastMessage},
			},
		},
	}
	_, err := m.collection.UpdateOne(context.TODO(), query, updateQuery)
	if err != nil {
		return &models.AppError{Code: models.DatabaseError, Message: err.Error()}
	}
	return nil
}

// AddJoinDate looks for the document with the userID provided.
// If found updates its server.JoinDates value.
// If not found an error is returned.
func (m *MongoDB) AddJoinDate(userID string, date string) *models.AppError {
	query := bson.M{
		"userID": userID,
	}
	var found *mongo.SingleResult
	found = m.collection.FindOne(context.TODO(), query)
	if found.Err() != nil && found.Err() == mongo.ErrNoDocuments {
		return &models.AppError{Code: models.UserNotFoundCode, Message: "User not found"}
	} else if found.Err() != nil {
		return &models.AppError{Code: models.DatabaseError, Message: found.Err().Error()}
	}
	updateQuery := bson.D{
		{
			Key: "$push", Value: bson.D{
				{Key: "server.joinDates", Value: date},
			},
		},
	}
	_, err := m.collection.UpdateOne(context.TODO(), query, updateQuery)
	if err != nil {
		return &models.AppError{Code: models.DatabaseError, Message: err.Error()}
	}
	return nil
}

// AddLeaveDate looks for the document with the userID provided.
// If found updates its server.leftDates value.
// If not found an error is returned.
func (m *MongoDB) AddLeaveDate(userID string, date string) *models.AppError {
	query := bson.M{
		"userID": userID,
	}
	var found *mongo.SingleResult
	found = m.collection.FindOne(context.TODO(), query)
	if found.Err() != nil && found.Err() == mongo.ErrNoDocuments {
		return &models.AppError{Code: models.UserNotFoundCode, Message: "User not found"}
	} else if found.Err() != nil {
		return &models.AppError{Code: models.DatabaseError, Message: found.Err().Error()}
	}
	updateQuery := bson.D{
		{
			Key: "$push", Value: bson.D{
				{Key: "server.leftDates", Value: date},
			},
		},
	}
	_, err := m.collection.UpdateOne(context.TODO(), query, updateQuery)
	if err != nil {
		return &models.AppError{Code: models.DatabaseError, Message: err.Error()}
	}
	return nil
}
