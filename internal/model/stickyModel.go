package model

import (
	"go.mongodb.org/mongo-driver/v2/bson"
)

type StickyChannel struct {
	ID bson.ObjectID `bson:"_id,omitempty"`
	ChannelID string `bson:"channelId,omitempty"`
	RecentPostMessageId string `bson:"recentPostMessageId,omitempty"`
	LastStickyMessageId string `bson:"lastStickyMessageId,omitempty"`
	StickyMessage string `bson:"stickyMessage,omitempty"`
}