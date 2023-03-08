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

func (r *Repository) catalogQueryFirstLine(ctx context.Context, contentLimit int64, filter bson.D) (*mongo.Cursor, error) {
	opt := &options.FindOptions{}
	opt.SetLimit(contentLimit)

	return r.Catalog().Find(ctx, filter, opt)
}
func (r *Repository) CatalogFirstLine(ctx context.Context, contentLimit int64) ([]DocCatalog, error) {
	curr, err := r.catalogQueryFirstLine(ctx, contentLimit, bson.D{})
	if err != nil {
		return nil, err
	}

	return DecodeCatalogCursor(ctx, curr)
}

func (r *Repository) CatalogFirstPageWithCategory(ctx context.Context, contentLimit int64, category string) ([]DocCatalog, error) {
	filter := bson.D{{Key: "kategori", Value: category}}

	curr, err := r.catalogQueryFirstLine(ctx, contentLimit, filter)
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

	return r.Catalog().Find(ctx, filter, opt)
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

// get catalog from greater than @id
func (r *Repository) CatalogGtId(ctx context.Context, id primitive.ObjectID, contentLimit int64) ([]DocCatalog, error) {
	gtID := bson.D{{Key: "$gt", Value: id}}
	return r.catalogIdSelector(ctx, gtID, contentLimit)
}

// get catalog from greater or equal @id
func (r *Repository) CatalogGteId(ctx context.Context, id primitive.ObjectID, contentLimit int64) ([]DocCatalog, error) {
	gtID := bson.D{{Key: "$gte", Value: id}}
	return r.catalogIdSelector(ctx, gtID, contentLimit)
}

func (r *Repository) Insert(ctx context.Context, in DocCatalog) (interface{}, error) {
	doc := DocumentFromModel(in)
	res, err := r.Catalog().InsertOne(ctx, doc)
	return res.InsertedID, err
}

func (r *Repository) CatalogFindOne(ctx context.Context, filter bson.D) (DocCatalog, error) {
	var doc DocCatalog
	err := r.Catalog().FindOne(context.TODO(), filter).Decode(&doc)

	return doc, err
}

func (r *Repository) FirstItem() (DocCatalog, error) {
	var doc DocCatalog
	findOptions := options.FindOne()
	// Sort by id field descending
	findOptions.SetSort(bson.D{{Key: "_id", Value: 1}})
	err := r.Catalog().FindOne(context.TODO(), bson.D{}).Decode(&doc)

	return doc, err
}

func (r *Repository) LastItem() (DocCatalog, error) {
	findOptions := options.FindOne()
	// Sort by id field descending
	findOptions.SetSort(bson.D{{Key: "_id", Value: -1}})

	var doc DocCatalog
	err := r.Catalog().FindOne(context.TODO(), bson.D{}, findOptions).Decode(&doc)

	return doc, err
}

func (r *Repository) CountCatalog(ctx context.Context) (int64, error) {
	filter := bson.D{}
	return r.Catalog().CountDocuments(ctx, filter)
}
