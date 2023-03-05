package repository

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Repository struct {
	Colletion
}

func NewRepository(db *mongo.Database) Repository {
	cimp := collectionImp{db: db}
	return Repository{Colletion: &cimp}
}

// OPTIONAL TODO
// projection filter field you want to return
//
//	opt.SetProjection(bson.D{
//		{Key: "field1", Value: 0,},
//		{Key: "field2", Value: 1},
//	})
func (r *Repository) CatalogFirstPage(ctx context.Context, contentLimit int64) ([]DocCatalog, error) {
	opt := &options.FindOptions{}
	opt.SetLimit(contentLimit)

	curr, err := r.Catalog().Find(ctx, nil, opt)
	if err != nil {
		return nil, err
	}

	return DecodeCatalogCursor(ctx, curr)
}

func (r *Repository) CatalogPageFromId(ctx context.Context, id primitive.ObjectID, contentLimit int64) ([]DocCatalog, error) {
	opt := &options.FindOptions{}
	opt.SetLimit(contentLimit)

	filter := bson.D{
		{Key: "_id", Value: bson.D{
			{Key: "$gt", Value: id},
		}},
	}

	curr, cErr := r.Catalog().Find(ctx, filter, opt)
	if cErr != nil {
		return nil, cErr
	}

	return DecodeCatalogCursor(ctx, curr)
}

func (r *Repository) CountCatalog(ctx context.Context) (int64, error) {
	return r.Catalog().CountDocuments(ctx, nil)
}
