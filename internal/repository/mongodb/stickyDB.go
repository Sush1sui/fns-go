package mongodb

import (
	"context"
	"fmt"

	"github.com/Sush1sui/fns-go/internal/common"
	"github.com/Sush1sui/fns-go/internal/model"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type MongoClient struct {
	Client *mongo.Collection
}

func (c MongoClient) InitializeStickyChannels() error {
	res, err := c.GetAllStickyChannels()
	if err != nil {
		return err
	}

	for _, channel := range res {
		common.StickyChannels[channel.ChannelID] = struct{}{} // Add the new channel ID to the StickyChannels map
	}
	return nil
}

func (c MongoClient) GetAllStickyChannels() ([]*model.StickyChannel, error) {
	var channels []*model.StickyChannel // Initialize an empty slice to hold the channels

	// Use the Find method to retrieve all sticky channels from the MongoDB collection
	cursor, err := c.Client.Find(context.Background(), bson.M{})
	if err != nil {
		return nil, fmt.Errorf("error finding sticky channels: %v", err)
	}
	defer cursor.Close(context.Background()) // Ensure the cursor is closed after use

	// Iterate through the cursor to decode each StickyChannel document into the channels slice
	for cursor.Next(context.Background()) {
		var channel model.StickyChannel // Create a new StickyChannel instance for each document
		if err := cursor.Decode(&channel); err != nil {
			return nil, fmt.Errorf("error decoding sticky channel: %v", err)
		}
		channels = append(channels, &channel) // Append the decoded channel to the slice
	}
	return channels, nil // Return the slice of StickyChannel pointers
}

func (c MongoClient) GetStickyChannel(channelID string) *model.StickyChannel {
	if channelID == "" {
		return nil // Return nil if the channelID is empty
	}

	var channel model.StickyChannel

	// Find the StickyChannel by ChannelID
	err := c.Client.FindOne(context.Background(), bson.M{"channel_id": channelID}).Decode(&channel)
	if err != nil {
		return nil // Return nil if no channel is found or if an error occurs
	}

	return &channel // Return the found StickyChannel
}

func (c MongoClient) CreateStickyChannel(channelID, stickyMessage string) (string, error) {
	if channelID == "" {
		return "", fmt.Errorf("channelID cannot be empty")
	}
	if stickyMessage == "" {
		stickyMessage = "Kindly avoid chatting or flood replies. Just use the **Thread** to avoid spamming or you will be **Timed out**"
	}

	channel := &model.StickyChannel{
		ChannelID:        channelID,
		StickyMessage:    stickyMessage,
		RecentPostMessageId: "",
		LastStickyMessageId: "",
	}

	// Insert the new StickyChannel into the MongoDB collection
	result, err := c.Client.InsertOne(context.Background(), channel)
	if err != nil {
		return "", fmt.Errorf("error inserting sticky channel: %v", err)
	}

	common.StickyChannels[channelID] = struct{}{} // Add the new channel ID to the StickyChannels map

	return result.InsertedID.(bson.ObjectID).Hex(), nil // Return the ID of the created StickyChannel
}

func (c MongoClient) GetRecentPostMessageId(channelId string) (string, error) {
	if channelId == "" {
		return "", fmt.Errorf("channelId cannot be empty")
	}
	var channel model.StickyChannel

	// Find the StickyChannel by ChannelID
	err := c.Client.FindOne(context.Background(), bson.M{"channel_id": channelId}).Decode(&channel)
	if err != nil {
		return "", fmt.Errorf("error finding sticky channel for channelId: %s, %v", channelId, err)
	}

	if channel.RecentPostMessageId == "" {
		return "", fmt.Errorf("no recent post message ID found for channelId: %s", channelId)
	}

	return channel.RecentPostMessageId, nil // Return the RecentPostMessageId
}

func (c MongoClient) UpdateStickyMessageId(channelId, lastMessageId, recentPostMessageId string) (*model.StickyChannel, error) {
	if channelId == "" {
		return nil, fmt.Errorf("channelId cannot be empty")
	}

	update := bson.M{
		"$set": bson.M{
			"last_sticky_message_id": lastMessageId,
			"recent_post_message_id": recentPostMessageId,
		},
	}

	result := c.Client.FindOneAndUpdate(context.Background(), bson.M{"channel_id": channelId}, update)
	if result.Err() != nil {
		return nil, fmt.Errorf("error updating sticky message ID for channelId: %s, %v", channelId, result.Err())
	}

	var updatedChannel model.StickyChannel
	if err := result.Decode(&updatedChannel); err != nil {
		return nil, fmt.Errorf("error decoding updated sticky channel: %v", err)
	}

	return &updatedChannel, nil
}

func (c MongoClient) SetStickyMessageId(channelId, stickyMessageId string) (*model.StickyChannel, error) {
	if channelId == "" {
		return nil, fmt.Errorf("channelId cannot be empty")
	}

	update := bson.M{
		"$set": bson.M{
			"sticky_message": stickyMessageId,
		},
	}

	result := c.Client.FindOneAndUpdate(context.Background(), bson.M{"channel_id": channelId}, update)
	if result.Err() != nil {
		return nil, fmt.Errorf("error updating sticky message ID for channelId: %s, %v", channelId, result.Err())
	}

	var updatedChannel model.StickyChannel
	if err := result.Decode(&updatedChannel); err != nil {
		return nil, fmt.Errorf("error decoding updated sticky channel: %v", err)
	}

	return &updatedChannel, nil
}

func (c MongoClient) DeleteStickyChannel(channelId string) (int, error) {
	if channelId == "" {
		return 0, fmt.Errorf("channelId cannot be empty")
	}

	result, err := c.Client.DeleteOne(context.Background(), bson.M{"channel_id": channelId})
	if err != nil {
		return 0, fmt.Errorf("error deleting sticky channel with channelId: %s, %v", channelId, err)
	}

	if result.DeletedCount == 0 {
		return 0, fmt.Errorf("no sticky channel found with channelId: %s", channelId)
	}

	// Remove the channel ID from the StickyChannels map
	delete(common.StickyChannels, channelId)

	return int(result.DeletedCount), nil // Return the number of deleted channels
}

func (c MongoClient) DeleteAllStickyChannels() (int, error) {
	result, err := c.Client.DeleteMany(context.Background(), bson.M{})
	if err != nil {
		return 0, fmt.Errorf("error deleting all sticky channels: %v", err)
	}

	// Clear the StickyChannels map after deletion
	common.StickyChannels = map[string]struct{}{}

	return int(result.DeletedCount), nil // Return the number of deleted channels
}