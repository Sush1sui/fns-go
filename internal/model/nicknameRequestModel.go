package model

import "go.mongodb.org/mongo-driver/v2/bson"

type NicknameRequest struct {
	ID bson.ObjectID `bson:"_id,omitempty"`
	UserID string `bson:"userId"`
	UserMessageID string `bson:"userMessageId"`
	UserChannelID string `bson:"userChannelId"`
	StaffChannelID string `bson:"staffChannelId"`
	StaffMessageID string `bson:"staffMessageId"`
	Nickname string `bson:"nickname"`
	Reactions []struct {
		Emoji string `bson:"emoji"`
	} `bson:"reactions"`
}