package repository

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Repository struct {
	coll Colletion
}

func NewRepository(db *mongo.Database) Repository {
	cimp := collectionImp{db: db}
	return Repository{coll: &cimp}
}

// OPTIONAL TODO
// projection filter field you want to return
//
//	opt.SetProjection(bson.D{
//		{Key: "field1", Value: 0,},
//		{Key: "field2", Value: 1},
//	})

func (r *Repository) catalogQueryFirstPage(ctx context.Context, contentLimit int64, filter bson.D) (*mongo.Cursor, error) {
	opt := &options.FindOptions{}
	opt.SetLimit(contentLimit)

	return r.coll.Catalog().Find(ctx, filter, opt)
}
func (r *Repository) CatalogFirstPage(ctx context.Context, contentLimit int64) ([]DocCatalog, error) {
	curr, err := r.catalogQueryFirstPage(ctx, contentLimit, bson.D{})
	if err != nil {
		return nil, err
	}

	return DecodeCatalogCursor(ctx, curr)
}

func (r *Repository) CatalogFirstPageWithCategory(ctx context.Context, contentLimit int64, category string) ([]DocCatalog, error) {
	filter := bson.D{{Key: "kategori", Value: category}}

	curr, err := r.catalogQueryFirstPage(ctx, contentLimit, filter)
	if err != nil {
		return nil, err
	}

	return DecodeCatalogCursor(ctx, curr)
}

// param @id primitive.ObjectID or set Query Selectors
func (r *Repository) catalogFromId(ctx context.Context, id interface{}, contentLimit int64) (*mongo.Cursor, error) {
	opt := &options.FindOptions{}
	opt.SetLimit(contentLimit)
	filter := bson.D{{Key: "_id", Value: id}}

	return r.coll.Catalog().Find(ctx, filter, opt)
}

// search the catalog with the underlying query selector
// ex Query Selectors bson.D{{Key: "$gt", Value: id}}
func (r *Repository) catalogIdSelector(ctx context.Context, id bson.D, contentLimit int64) ([]DocCatalog, error) {
	curr, err := r.catalogFromId(ctx, id, contentLimit)
	if err != nil {
		return nil, err
	}
	return DecodeCatalogCursor(ctx, curr)
}

func (r *Repository) CataloGtId(ctx context.Context, id primitive.ObjectID, contentLimit int64) ([]DocCatalog, error) {
	gtID := bson.D{{Key: "$gt", Value: id}}
	return r.catalogIdSelector(ctx, gtID, contentLimit)
}

func (r *Repository) CatalogGteId(ctx context.Context, id primitive.ObjectID, contentLimit int64) ([]DocCatalog, error) {
	gtID := bson.D{{Key: "$gte", Value: id}}
	return r.catalogIdSelector(ctx, gtID, contentLimit)
}

func (r *Repository) LastItem() (DocCatalog, error) {
	findOptions := options.FindOne()
	// Sort by id field descending
	findOptions.SetSort(bson.D{{Key: "_id", Value: -1}})

	var doc DocCatalog
	err := r.coll.Catalog().FindOne(context.TODO(), bson.D{}, findOptions).Decode(&doc)

	return doc, err
}

func (r *Repository) FirstItem() (DocCatalog, error) {
	var doc DocCatalog
	err := r.coll.Catalog().FindOne(context.TODO(), bson.D{}).Decode(&doc)

	return doc, err
}

func (r *Repository) CountCatalog(ctx context.Context) (int64, error) {
	filter := bson.D{}
	return r.coll.Catalog().CountDocuments(ctx, filter)
}
