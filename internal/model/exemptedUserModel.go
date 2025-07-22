package model

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type ExemptedUser struct {
	ID          bson.ObjectID `bson:"_id,omitempty"`
	UserID      string         `bson:"userId,omitempty"`
	Expiration  time.Time     `bson:"expiration,omitempty"`
}