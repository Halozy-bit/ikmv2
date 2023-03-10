package sidejob

import (
	"context"
	"log"
	"time"

	asynctask "github.com/ikmv2/backend/pkg/async_task"
	"github.com/ikmv2/backend/pkg/cache"
	"github.com/ikmv2/backend/pkg/helper"
	"github.com/ikmv2/backend/pkg/repository"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Task 2
type RefreshCatalogCategoryPage struct {
	asynctask.TaskIdentifier
	Db *mongo.Database
}

// TODO
// get all category
// set category to slice
// loop all category
// nest loop in category per page then save to cache
func (rcp RefreshCatalogCategoryPage) Run() {
	coll := rcp.Db.Collection("catalog")

	opt := options.Find()
	opt.SetProjection(bson.D{{Key: "_id", Value: 1}})

	// NOTE
	// helper.CategoryAvail containing list of string
	// val is string
	for _, val := range helper.CategoryAvail {
		filter := bson.D{{Key: repository.CategoryField, Value: val}}
		ctlgTotal, err := coll.CountDocuments(context.TODO(), filter)
		if err != nil || ctlgTotal == 0 {
			log.Println("no product to refresh")
			return
		}

		err = rotate_Ctgry_PerPage(coll, val, int(ctlgTotal), opt)
		if err != nil {
			log.Println(err)
			return
		}
	}
}

func rotate_Ctgry_PerPage(coll *mongo.Collection, category string, itemTotal int, opt *options.FindOptions) error {
	maxPage := helper.MaxPage(helper.MaxProductPerPage, itemTotal)

	idTable := make([]primitive.ObjectID, maxPage)
	cache_LastID := cache.Pagination.CategoryPage(category, 1)
	last_id := initLastIDCategory(cache_LastID, coll, category, opt)

	for i := 0; i < maxPage; i++ {
		page := i + 1
		time.Sleep(time.Millisecond * 10)
		filter := bson.D{{Key: repository.CategoryField, Value: category}}

		if last_id != primitive.NilObjectID {
			id := bson.D{{Key: "$gt", Value: last_id}}
			filter = append(filter, bson.E{Key: "_id", Value: id})
		}

		totalNext := helper.CountTtlProductNxtPage(page, itemTotal)
		// findFirstAndLast is compatible with category because it defines its own filter
		fl, err := findFirstAndLast(coll, filter, totalNext, opt)
		if err != nil {
			return err
		}

		idTable[i] = fl[0]
		last_id = fl[1]
	}

	return cache.Pagination.StoreCategory(category, idTable)
}
