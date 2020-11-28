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
func (m *MongoDB) InitDB(c config.DBConnection) *dBError {
	var err error
	options := options.Client().ApplyURI(fmt.Sprintf("mongodb+srv://%s:%s@%s/%s%s", c.User, c.Pass, c.Host, c.DB, c.Options))
	m.client, err = mongo.Connect(context.TODO(), options)
	if err != nil {
		return &dBError{Code: CantConnectToDatabaseCode, Message: err.Error()}
	}
	err = m.client.Ping(context.TODO(), nil)
	if err != nil {
		return &dBError{Code: CantConnectToDatabaseCode, Message: err.Error()}
	}
	m.collection = m.client.Database(c.DB).Collection(c.Collection)
	return nil
}

// GetUser searches and returns user information from the database that matches either the userId or username.
// It can accept either username or userId or both, but if both are empty an error will be returned.
func (m *MongoDB) GetUser(userID string, username string) (models.User, *dBError) {
	var query bson.M
	if userID != "" && username != "" {
		query = bson.M{"$or": []bson.M{{"userID": userID}, {"fullName": username}}}
	} else if userID != "" && username == "" {
		query = bson.M{"userID": userID}
	} else if userID == "" && username != "" {
		query = bson.M{"username": username}
	} else {
		return models.User{}, &dBError{Code: WrongParametersCode, Message: WrongParametersMessage}
	}
	var result models.User
	res := m.collection.FindOne(context.TODO(), query)
	if res.Err() != nil && res.Err() == mongo.ErrNoDocuments {
		return models.User{}, &dBError{Code: UserNotFoundCode, Message: UserNotFoundMessage}
	} else if res.Err() != nil {
		return models.User{}, &dBError{Code: DatabaseErrorCode, Message: res.Err().Error()}
	}
	err := res.Decode(&result)
	if err != nil {
		return models.User{}, &dBError{Code: DecodingErrorCode, Message: err.Error()}
	}
	return result, nil
}

// AddUser adds the user information to the database.
func (m *MongoDB) AddUser(user models.User) *dBError {
	_, gErr := m.GetUser(user.UserID, "")
	if gErr == nil {
		return &dBError{Code: UserAlredyExistsCode, Message: UserAlredyExistsMessage}
	} else if gErr.Code != UserNotFoundCode {
		return gErr
	}
	_, err := m.collection.InsertOne(context.TODO(), user)
	if err != nil {
		return &dBError{Code: DatabaseErrorCode, Message: err.Error()}
	}
	return nil
}

// IncreaseMessageCount searches for the document with the provided userID.
// If found, it increases its server.messageCount value and updates server.lastMessage.
// If the user is not found an error is returned.
func (m *MongoDB) IncreaseMessageCount(userID string) *dBError {
	query := bson.M{
		"userID": userID,
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
	ur, err := m.collection.UpdateOne(context.TODO(), query, updateQuery)
	if err != nil {
		return &dBError{Code: DatabaseErrorCode, Message: err.Error()}
	}
	if ur.MatchedCount == 0 {
		return &dBError{Code: UserNotFoundCode, Message: UserNotFoundMessage}
	}
	return nil
}

// AddJoinDate looks for the document with the userID provided.
// If found updates its server.JoinDates value.
// If not found an error is returned.
func (m *MongoDB) AddJoinDate(userID string, date time.Time) (bool, *dBError) {
	user, dbErr := m.GetUser(userID, "")
	if dbErr != nil {
		return false, dbErr
	}
	query := bson.M{
		"userID": userID,
	}
	date = date.In(config.Loc)
	updateQuery := bson.D{
		{
			Key: "$push", Value: bson.D{
				{Key: "server.joinDates", Value: date.Format(time.RFC822)},
			},
		},
	}
	_, err := m.collection.UpdateOne(context.TODO(), query, updateQuery)
	if err != nil {
		return false, &dBError{Code: DatabaseErrorCode, Message: err.Error()}
	}
	return user.Server.Ultimatum, nil
}

// AddLeaveDate looks for the document with the userID provided.
// If found updates its server.leftDates value.
// If not found an error is returned.
func (m *MongoDB) AddLeaveDate(userID string, date time.Time) (bool, *dBError) {
	user, dbErr := m.GetUser(userID, "")
	if dbErr != nil {
		return false, dbErr
	}
	findQuery := bson.M{
		"userID": userID,
	}
	date.In(config.Loc)
	updateQuery := bson.D{
		{
			Key: "$push", Value: bson.D{
				{Key: "server.leftDates", Value: date.Format(time.RFC822)},
			},
		},
	}
	_, err := m.collection.UpdateOne(context.TODO(), findQuery, updateQuery)
	if err != nil {
		return false, &dBError{Code: DatabaseErrorCode, Message: err.Error()}
	}
	return user.Server.Ultimatum, nil
}

func (m *MongoDB) SetUltimatum(userID string) *dBError {
	user, dbErr := m.GetUser(userID, "")
	if dbErr != nil {
		return dbErr
	}
	if user.Server.Ultimatum {
		return &dBError{Code: UserAlredyInUltimatumCode, Message: UserAlredyInUltimatumMessage}
	}
	filter := bson.M{
		"userID": userID,
	}
	updateQuery := bson.D{
		{
			Key: "$set", Value: bson.D{
				{Key: "server.ultimatum", Value: true},
			},
		},
	}
	_, err := m.collection.UpdateOne(context.TODO(), filter, updateQuery)
	if err != nil {
		return &dBError{Code: DatabaseErrorCode, Message: err.Error()}
	}
	return nil
}
