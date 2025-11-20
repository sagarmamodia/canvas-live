package model

import "go.mongodb.org/mongo-driver/bson/primitive"

type Slide struct {
	ID         string   `bson:"_id" json:"id"`
	Background string   `bson:"background" json:"background"`
	Objects    []Object `bson:"objects" json:"objects"`
}

type Document struct {
	ID      primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Title   string             `bson:"title" json:"title"`
	OwnerID string             `bson:"ownerId" json:"ownerId"`
	Slides  []Slide            `bson:"slides" json:"slides"`
}

type Object struct {
	ID         string                 `bson:"_id" json:"id"`
	Type       string                 `bson:"type" json:"type"`
	Attributes map[string]interface{} `bson:"attributes" json:"attributes"`
}
