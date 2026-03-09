package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Card struct {
	ID                   primitive.ObjectID `bson:"_id,omitempty"`
	Code                 string             `bson:"code"`
	Name                 string             `bson:"name"`
	ArcanaType           string             `bson:"arcanaType"`
	Description          string             `bson:"description"`
	InterpretationOpen   string             `bson:"interpretationOpen"`
	InterpretationClosed string             `bson:"interpretationClosed"`
}
