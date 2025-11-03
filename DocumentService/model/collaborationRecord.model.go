package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CollaborationRecord struct {
	ID         primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	UserID     string             `bson:"userId" json:"userId"`
	DocumentID string             `bson:"documentId" json:"documentId"`
	AccessType string             `bson:"accessType" json:"accessType"` // {Editor, Viewer}
	SharedAt   time.Time          `bson:"sharedAt" json:"sharedAt"`
}
