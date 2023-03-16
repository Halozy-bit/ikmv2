package repository

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

const (
	CategoryField = "kategori"
)

// TODO
// Foto change to array focuse on per product
// Thumbnail for displaying catalog

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
	Weight      []string           `bson:"ukuran" json:"ukuran"`
	Variant     []string           `bson:"varian" json:"varian"`
	Composition []string           `bson:"komposisi" json:"komposisi"`
	Description string             `bson:"deskripsi" json:"deskripsi"`
	Owner       string             `bson:"owner" json:"owner"`
	Foto        Foto               `bson:"foto" json:"foto"`
}

type Foto struct {
	Cover   string `bson:"cover" json:"cover"`
	Detail1 string `bson:"detail1" json:"detail1"`
	Detail2 string `bson:"detail2" json:"detail2"`
}

func ProductToDocument(p Product) bson.D {
	return bson.D{
		{Key: "nama", Value: p.Name},
		{Key: "kategori", Value: p.Category},
		{Key: "ukuran", Value: p.Weight},
		{Key: "varian", Value: p.Variant},
		{Key: "komposisi", Value: p.Composition},
		{Key: "deskripsi", Value: p.Description},
		{Key: "owner", Value: p.Owner},
		{Key: "foto", Value: bson.D{
			{Key: "cover", Value: p.Foto.Cover},
			{Key: "detail1", Value: p.Foto.Detail1},
			{Key: "detail2", Value: p.Foto.Detail2},
		}},
	}
}
