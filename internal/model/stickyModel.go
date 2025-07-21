package model

import (
	"go.mongodb.org/mongo-driver/v2/bson"
)

type StickyChannel struct {
	ID bson.ObjectID `bson:"_id,omitempty"`
	ChannelID string `bson:"channel_id,omitempty"`
	RecentPostMessageId string `bson:"recent_post_message_id,omitempty"`
	LastStickyMessageId string `bson:"last_sticky_message_id,omitempty"`
	StickyMessage string `bson:"sticky_message,omitempty"`
}