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

func (r *Repository) CatalogFirstLine(ctx context.Context, contentLimit int64) ([]CatalogDisplay, error) {
	opt := options.Find().SetLimit(contentLimit)
	return r.catalogFinder(ctx, bson.D{}, opt)
}

func (r *Repository) CatalogFirstPageWithCategory(ctx context.Context, contentLimit int64, category string) ([]CatalogDisplay, error) {
	opt := options.Find().SetLimit(contentLimit)
	filter := bson.D{{Key: CategoryField, Value: category}}
	return r.catalogFinder(ctx, filter, opt)
}

func (r *Repository) catalogFinder(ctx context.Context, filter bson.D, opt *options.FindOptions) ([]CatalogDisplay, error) {
	opt.SetProjection(bson.D{
		{Key: "_id", Value: 1}, {Key: "nama", Value: 1}, {Key: "kategori", Value: 1},
		{Key: "owner", Value: 1}, {Key: "thumbnail", Value: "$foto.cover"},
	})

	curr, err := r.Catalog().Find(ctx, filter, opt)
	if err != nil {
		return nil, err
	}
	return DecodeCatalogCursor(ctx, curr)
}

// search the catalog with the underlying query selector
// ex Query Selectors bson.D{{Key: "$gt", Value: id}}
func (r *Repository) CatalogIdSelector(ctx context.Context, contentLimit int64, id bson.D, addFilter ...bson.E) ([]CatalogDisplay, error) {
	opt := &options.FindOptions{}
	opt.SetLimit(contentLimit)

	filter := bson.D{{Key: "_id", Value: id}}
	if len(addFilter) > 0 {
		filter = append(filter, addFilter...)
	}

	return r.catalogFinder(ctx, filter, opt)
}

// get catalog from greater than @id
func (r *Repository) CatalogGtId(ctx context.Context, id primitive.ObjectID, contentLimit int64, category ...string) ([]CatalogDisplay, error) {
	gtID := bson.D{{Key: "$gt", Value: id}}
	if len(category) > 0 {
		return r.CatalogIdSelector(ctx, contentLimit, gtID, bson.E{Key: CategoryField, Value: category})
	}
	return r.CatalogIdSelector(ctx, contentLimit, gtID)
}

// get catalog from greater or equal @id
func (r *Repository) CatalogGteId(ctx context.Context, id primitive.ObjectID, contentLimit int64, category ...string) ([]CatalogDisplay, error) {
	gtID := bson.D{{Key: "$gte", Value: id}}
	return r.CatalogIdSelector(ctx, contentLimit, gtID)
}

func (r *Repository) CatalogGteIdByCategory(ctx context.Context, id primitive.ObjectID, contentLimit int64, category string) ([]CatalogDisplay, error) {
	gtID := bson.D{{Key: "$gte", Value: id}}
	return r.CatalogIdSelector(ctx, contentLimit, gtID, bson.E{Key: CategoryField, Value: category})
}

func (r *Repository) InsertCatalog(ctx context.Context, doc bson.D) (interface{}, error) {
	res, err := r.Catalog().InsertOne(ctx, doc)
	return res.InsertedID, err
}

func (r *Repository) FirstItem() (CatalogDisplay, error) {
	var doc CatalogDisplay
	findOptions := options.FindOne()
	// Sort by id field ascending
	findOptions.SetSort(bson.D{{Key: "_id", Value: 1}})
	err := r.Catalog().FindOne(context.TODO(), bson.D{}).Decode(&doc)

	return doc, err
}

func (r *Repository) LastItem() (CatalogDisplay, error) {
	findOptions := options.FindOne()
	// Sort by id field descending
	findOptions.SetSort(bson.D{{Key: "_id", Value: -1}})

	var doc CatalogDisplay
	err := r.Catalog().FindOne(context.TODO(), bson.D{}, findOptions).Decode(&doc)

	return doc, err
}

func (r *Repository) CountCatalog(ctx context.Context) (int64, error) {
	filter := bson.D{}
	return r.Catalog().CountDocuments(ctx, filter)
}

func (r *Repository) FindProduct(ctx context.Context, id primitive.ObjectID) (Product, error) {
	filter := bson.D{{Key: "_id", Value: id}}

	var p Product
	err := r.Catalog().FindOne(ctx, filter).Decode(&p)
	return p, err
}

func (r *Repository) CountCatalogCategory(ctx context.Context, category string) (int64, error) {
	filter := bson.D{bson.E{Key: CategoryField, Value: category}}
	return r.Catalog().CountDocuments(ctx, filter)
}
