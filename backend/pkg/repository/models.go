package repository

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

const (
	CategoryField  = "kategori"
	ThumbnailField = "thumbnail"
)

// TODO
// Foto change to array focuse on per product
// Thumbnail for displaying catalog

type DocCatalog struct {
	Id          primitive.ObjectID `bson:"_id" insert:"0"`
	Name        string             `bson:"nama"`
	Category    string             `bson:"kategori"`
	Description string             `bson:"deskripsi"`
	Owner       string             `bson:"owner"`
	Thumbnail   string             `bson:"thumbnail,omitempty"`
	Foto        string             `bson:"foto"`
}

type CatalogDisplay struct {
	Id        primitive.ObjectID `bson:"_id" json:"id"`
	Name      string             `bson:"nama" json:"nama"`
	Category  string             `bson:"kategori" json:"kategori"`
	Owner     string             `bson:"owner" json:"owner"`
	Thumbnail string             `bson:"thumbnail" json:"thumbnail"`
}

type Product struct {
	Id          primitive.ObjectID `bson:"_id" json:"id"`
	Name        string             `bson:"nama" json:"nama"`
	Category    string             `bson:"kategori" json:"kategori"`
	Description string             `bson:"deskripsi" json:"deskripsi"`
	Owner       string             `bson:"owner" json:"owner"`
	Foto        string             `bson:"foto" json:"foto"`
}
