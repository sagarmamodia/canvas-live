package model

import "go.mongodb.org/mongo-driver/bson/primitive"

type Slide struct {
	ID         string        `bson:"_id,omitemtpy" json:"id,omitempty"`
	Background string        `bson:"background" json:"background"`
	Objects    []interface{} `bson:"objects" json:"objects"`
}

type Document struct {
	ID      primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Title   string             `bson:"title" json:"title"`
	OwnerID string             `bson:"ownerId" json:"ownerId"`
	Slides  []Slide            `bson:"slides" json:"slides"`
}
