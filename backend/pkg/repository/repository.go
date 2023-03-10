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

func (r *Repository) CatalogFirstLine(ctx context.Context, contentLimit int64) ([]DocCatalog, error) {
	opt := options.Find().SetLimit(contentLimit)
	return r.catalogFinder(ctx, bson.D{}, opt)
}

func (r *Repository) CatalogFirstPageWithCategory(ctx context.Context, contentLimit int64, category string) ([]DocCatalog, error) {
	opt := options.Find().SetLimit(contentLimit)
	filter := bson.D{{Key: CategoryField, Value: category}}
	return r.catalogFinder(ctx, filter, opt)
}

func (r *Repository) catalogFinder(ctx context.Context, filter bson.D, opt *options.FindOptions) ([]DocCatalog, error) {
	curr, err := r.Catalog().Find(ctx, filter, opt)
	if err != nil {
		return nil, err
	}
	return DecodeCatalogCursor(ctx, curr)
}

// search the catalog with the underlying query selector
// ex Query Selectors bson.D{{Key: "$gt", Value: id}}
func (r *Repository) CatalogIdSelector(ctx context.Context, contentLimit int64, id bson.D, addFilter ...bson.E) ([]DocCatalog, error) {
	opt := &options.FindOptions{}
	opt.SetLimit(contentLimit)

	filter := bson.D{{Key: "_id", Value: id}}
	if len(addFilter) > 0 {
		filter = append(filter, addFilter...)
	}

	return r.catalogFinder(ctx, filter, opt)
}

// get catalog from greater than @id
func (r *Repository) CatalogGtId(ctx context.Context, id primitive.ObjectID, contentLimit int64, category ...string) ([]DocCatalog, error) {
	gtID := bson.D{{Key: "$gt", Value: id}}
	if len(category) > 0 {
		return r.CatalogIdSelector(ctx, contentLimit, gtID, bson.E{Key: CategoryField, Value: category})
	}
	return r.CatalogIdSelector(ctx, contentLimit, gtID)
}

// get catalog from greater or equal @id
func (r *Repository) CatalogGteId(ctx context.Context, id primitive.ObjectID, contentLimit int64, category ...string) ([]DocCatalog, error) {
	gtID := bson.D{{Key: "$gte", Value: id}}
	return r.CatalogIdSelector(ctx, contentLimit, gtID)
}

func (r *Repository) CatalogGteIdByCategory(ctx context.Context, id primitive.ObjectID, contentLimit int64, category string) ([]DocCatalog, error) {
	gtID := bson.D{{Key: "$gte", Value: id}}
	return r.CatalogIdSelector(ctx, contentLimit, gtID, bson.E{Key: CategoryField, Value: category})
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
	// Sort by id field ascending
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

func (r *Repository) CountCatalogCategory(ctx context.Context, category string) (int64, error) {
	filter := bson.D{bson.E{Key: CategoryField, Value: category}}
	return r.Catalog().CountDocuments(ctx, filter)
}
