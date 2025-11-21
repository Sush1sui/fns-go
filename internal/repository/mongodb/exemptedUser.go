package mongodb

import (
	"context"
	"fmt"
	"time"

	"github.com/Sush1sui/fns-go/internal/model"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

func (c *MongoClient) ExemptUserVanity(userID, role string) (bool, error) {
	if userID == "" {
		return false, fmt.Errorf("userID cannot be empty")
	}
	
	var res bson.M
	if role == "staff" || userID == "1258348384671109120" {
		err := c.Client.FindOneAndUpdate(
			context.Background(),
			bson.M{"userId": userID}, // filter
			bson.M{ // update
				"$set": bson.M{
                    "userId": userID,
                },
                "$unset": bson.M{
                    "expiration": "",
                },
			},
			options.FindOneAndUpdate().SetUpsert(true).SetReturnDocument(options.After),
		).Decode(&res)

		if err != nil {
			fmt.Println("Error updating or inserting exempted user:", err)
			return false, err
		}

		return true, nil
	}

	expirationDate := time.Now().Add(3 * 24 * time.Hour) // 3 days from now
	err := c.Client.FindOneAndUpdate(
		context.Background(),
		bson.M{"userId": userID}, // filter
		bson.M{ // update
			"$set": bson.M{
				"userId": userID,
				"expiration": expirationDate,
			},
		},
		options.FindOneAndUpdate().SetUpsert(true).SetReturnDocument(options.After),
	).Decode(&res)

	if err != nil {
		fmt.Println("Error updating or inserting exempted user:", err)
		return false, err
	}

	return true, nil
}

func (c *MongoClient) RemoveExemptedUser(userID string) (int, error) {
	if userID == "" {
		return 0, fmt.Errorf("userID cannot be empty")
	}

	res, err := c.Client.DeleteOne(context.Background(), bson.M{"userId": userID})
	if err != nil {
		return 0, fmt.Errorf("error deleting exempted user with userID: %s, %v", userID, err)
	}

	return int(res.DeletedCount), nil
}

func (c *MongoClient) GetAllExemptedUsers() ([]*model.ExemptedUser, error) {
	cursor, err := c.Client.Find(context.Background(), bson.M{})
	if err != nil {
		return nil, fmt.Errorf("error fetching exempted users: %v", err)
	}
	defer cursor.Close(context.Background())

	var exemptedUsers []*model.ExemptedUser
	for cursor.Next(context.Background()) {
		var user model.ExemptedUser
		if err := cursor.Decode(&user); err != nil {
			return nil, fmt.Errorf("error decoding exempted user: %v", err)
		}
		exemptedUsers = append(exemptedUsers, &user)
	}

	if err := cursor.Err(); err != nil {
		return nil, fmt.Errorf("cursor error: %v", err)
	}

	return exemptedUsers, nil
}