package repository

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type DocCatalog struct {
	Id          primitive.ObjectID `bson:"_id" insert:"0"`
	Name        string             `bson:"nama"`
	Category    string             `bson:"kategori"`
	Description string             `bson:"deskripsi"`
	Owner       string             `bson:"owner"`
	Foto        string             `bson:"foto"`
}
