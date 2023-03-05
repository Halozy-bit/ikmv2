package repository

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type DocCatalog struct {
	Id          primitive.ObjectID `bson:"_id"`
	Name        string             `bson:"name"`
	Category    []string           `bson:"category"`
	Description string             `bson:"description"`
	Owner       string             `bson:"owner"`
	Foto        []string           `bson:"foto"`
}
