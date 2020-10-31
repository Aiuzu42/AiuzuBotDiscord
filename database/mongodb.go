package database

import (
	"context"
	"fmt"

	"github.com/aiuzu42/AiuzuBotDiscord/config"
	"github.com/aiuzu42/AiuzuBotDiscord/models"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoDB struct {
	client     *mongo.Client
	collection *mongo.Collection
}

func (m *MongoDB) InitDB(c config.DBConnection) *models.AppError {
	var err error
	options := options.Client().ApplyURI(fmt.Sprintf("mongodb+srv://%s:%s@%s/%s%s", c.User, c.Pass, c.Host, c.DB, c.Options))
	m.client, err = mongo.Connect(context.TODO(), options)
	if err != nil {
		log.Error(err.Error())
		return &models.AppError{Code: models.CantConnectToDatabase, Message: "Error Connecting to database"}
	}
	err = m.client.Ping(context.TODO(), nil)
	if err != nil {
		log.Error(err.Error())
		return &models.AppError{Code: models.CantConnectToDatabase, Message: "Error Connecting to database"}
	}
	m.collection = m.client.Database(c.DB).Collection(c.Collection)
	return nil
}

func (m *MongoDB) GetUser(userID string, username string) (models.User, *models.AppError) {
	query := bson.M{
		"userID": userID,
		"$or":    bson.M{"fullName": username},
	}
	var result models.User
	err := m.collection.FindOne(context.TODO(), query).Decode(&result)
	if err != nil {
		return models.User{}, &models.AppError{Code: models.DatabaseError, Message: "Database Error"}
	}
	if result.UserID == "" {
		return models.User{}, &models.AppError{Code: models.UserNotFoundCode, Message: "User not found"}
	}
	return result, nil
}

func (m *MongoDB) AddUser(user models.User) *models.AppError {
	query := bson.M{
		"userID": user.UserID,
	}
	var found *mongo.SingleResult
	found = m.collection.FindOne(context.TODO(), query)
	if found.Err() == nil {
		return &models.AppError{Code: models.UserAlredyExists, Message: "The user alredy exists"}
	} else if found.Err() != nil && found.Err() != mongo.ErrNoDocuments {
		log.Error(found.Err().Error())
		return &models.AppError{Code: models.DatabaseError, Message: "Database Error"}
	}
	_, err := m.collection.InsertOne(context.TODO(), user)
	if err != nil {
		return &models.AppError{Code: models.DatabaseError, Message: "Database Error"}
	}
	return nil
}

func (m *MongoDB) IncreaseMessageCount(userID string) *models.AppError {
	query := bson.M{
		"userID": userID,
	}
	var found *mongo.SingleResult
	found = m.collection.FindOne(context.TODO(), query)
	if found.Err() != nil && found.Err() == mongo.ErrNoDocuments {
		return &models.AppError{Code: models.UserNotFoundCode, Message: "User not found"}
	} else if found.Err() != nil {
		log.Error(found.Err().Error())
		return &models.AppError{Code: models.DatabaseError, Message: "Database Error"}
	}
	updateQuery := bson.D{
		{
			Key: "$inc", Value: bson.D{
				{Key: "server.messageCount", Value: 1},
			},
		},
	}
	_, err := m.collection.UpdateOne(context.TODO(), query, updateQuery)
	if err != nil {
		return &models.AppError{Code: models.DatabaseError, Message: "Database Error"}
	}
	return nil
}

func (m *MongoDB) AddJoinDate(userID string, date string) *models.AppError {
	query := bson.M{
		"userID": userID,
	}
	var found *mongo.SingleResult
	found = m.collection.FindOne(context.TODO(), query)
	if found.Err() != nil && found.Err() == mongo.ErrNoDocuments {
		return &models.AppError{Code: models.UserNotFoundCode, Message: "User not found"}
	} else if found.Err() != nil {
		log.Error(found.Err().Error())
		return &models.AppError{Code: models.DatabaseError, Message: "Database Error"}
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
		return &models.AppError{Code: models.DatabaseError, Message: "Database Error"}
	}
	return nil
}
