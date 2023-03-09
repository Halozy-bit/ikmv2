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

type RefreshCatalogPage struct {
	asynctask.TaskIdentifier
	Db *mongo.Database
}

func (rc *RefreshCatalogPage) Run() {
	coll := rc.Db.Collection("catalog")
	ctlgTotal, err := coll.CountDocuments(context.TODO(), bson.D{})
	if err != nil || ctlgTotal == 0 {
		log.Println("no product to refresh")
		return
	}

	maxPage := helper.MaxPage(helper.MaxProductPerPage, int(ctlgTotal))

	opt := options.Find()
	opt.SetProjection(bson.D{
		{Key: "_id", Value: 1},
	})

	idTable := make([]primitive.ObjectID, maxPage)
	last_id := cache.Pagination.Page(1)
	if last_id != primitive.NilObjectID {
		// skip 6 document
		opt.SetLimit(5)
		filter := bson.D{{Key: "_id", Value: bson.D{{Key: "$gt", Value: last_id}}}}
		fl, err := findFirstAndLast(coll, filter, opt)
		if err != nil {
			log.Println(err)
			return
		}

		last_id = fl[1]
	}

	opt.SetLimit(int64(helper.MaxProductPerPage))
	for i := 0; i < maxPage; i++ {
		time.Sleep(time.Millisecond * 10)

		filter := bson.D{}
		if last_id != primitive.NilObjectID {
			id := bson.D{{Key: "$gt", Value: last_id}}
			filter = bson.D{{Key: "_id", Value: id}}
		}

		ids, err := findFirstAndLast(coll, filter, opt)
		if err != nil {
			log.Println(err)
			return
		}

		idTable[i] = ids[0]
		log.Println(ids[0])
		last_id = ids[1]
	}
	cache.Pagination.StorePage(idTable)
}

func findFirstAndLast(coll *mongo.Collection, filter bson.D, opt *options.FindOptions) ([2]primitive.ObjectID, error) {
	var ids [2]primitive.ObjectID
	curr, err := coll.Find(context.TODO(), filter, opt)
	if err != nil {
		return ids, err
	}

	ids, err = decodeFirstAndLastID(context.TODO(), curr)
	return ids, err
}

func decodeFirstAndLastID(ctx context.Context, curr *mongo.Cursor) ([2]primitive.ObjectID, error) {
	// var once *sync.Once
	type IDDecoder struct {
		Id primitive.ObjectID `bson:"_id"`
	}

	var result [2]primitive.ObjectID
	var i = 0
	for {
		var tmp IDDecoder
		var err error

		if !curr.TryNext(ctx) {
			err = curr.Decode(&tmp)
			if err != nil {
				return result, err
			}
			result[1] = tmp.Id
			break
		} else if i == 0 {
			err = curr.Decode(&tmp)
			if err != nil {
				return result, err
			}
			result[0] = tmp.Id
		}
		i++

	}

	return result, nil
}
