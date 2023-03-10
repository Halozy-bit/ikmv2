package sidejob

import (
	"context"
	"log"
	"time"

	asynctask "github.com/ikmv2/backend/pkg/async_task"
	"github.com/ikmv2/backend/pkg/cache"
	"github.com/ikmv2/backend/pkg/helper"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

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
	opt.SetLimit(int64(helper.MaxProductPerPage))

	for _, val := range helper.CategoryAvail {
		ctlgTotal, err := coll.CountDocuments(context.TODO(), bson.D{})
		if err != nil || ctlgTotal == 0 {
			log.Println("no product to refresh")
			return
		}

		rotate_Ctgry_PerPage(coll, val, int(ctlgTotal), opt)

	}
}

func rotate_Ctgry_PerPage(coll *mongo.Collection, category string, itemTotal int, opt *options.FindOptions) {
	maxPage := helper.MaxPage(helper.MaxProductPerPage, itemTotal)

	idTable := make([]primitive.ObjectID, maxPage)
	last_id := initLastID(cache.Pagination.CategoryPage(category, 1), coll, opt)

	for page := 0; page < maxPage; page++ {
		time.Sleep(time.Millisecond * 10)
		filter := bson.D{{Key: "category", Value: category}}

		if last_id != primitive.NilObjectID {
			id := bson.D{{Key: "$gt", Value: last_id}}
			filter = append(filter, bson.E{Key: "_id", Value: id})
		}

		totalNext := helper.CountTtlProductNxtPage(page, itemTotal)
		fl, err := findFirstAndLast(coll, filter, totalNext, opt)
		if err != nil {
			return
		}

		idTable[page] = fl[0]
		last_id = fl[1]
	}

	err := cache.Pagination.StoreCategory(category, idTable)
	if err != nil {
		log.Print(err)
	}
}
