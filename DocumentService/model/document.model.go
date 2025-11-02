package model

import "go.mongodb.org/mongo-driver/bson/primitive"

// DocumentData Model
type DocObject interface {
}

type Slide struct {
	ID         primitive.ObjectID `bson:"_id,omitemtpy" json:"id,omitempty"`
	Background string             `bson:"background" json:"background"`
	Objects    []DocObject        `bson:"objects" json:"objects"`
}

type Document struct {
	ID      primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Title   string             `bson:"title" json:"title"`
	OwnerID string             `bson:"ownerId" json:"ownerId"`
	Slides  []Slide            `bson:"slides" json:"slides"`
}
