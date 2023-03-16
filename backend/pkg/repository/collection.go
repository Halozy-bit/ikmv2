package repository

import "go.mongodb.org/mongo-driver/mongo"

type Colletion interface {
	Catalog() *mongo.Collection
}

type collectionImp struct {
	db *mongo.Database
}

func (c *collectionImp) Catalog() *mongo.Collection {
	return c.db.Collection("catalog")
}

func (c *collectionImp) Owner() *mongo.Collection {
	return c.db.Collection("owner")
}
