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

// Task 1
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
	opt.SetProjection(bson.D{{Key: "_id", Value: 1}})

	idTable := make([]primitive.ObjectID, maxPage)
	// WARN
	// Minor bug
	// last_id is the reference id to perform the next query
	last_id := initLastID(cache.Pagination.Page(1), coll, opt)

	for i := 0; i < maxPage; i++ {
		// start with the first page
		page := i + 1
		time.Sleep(time.Millisecond * 10)

		// total item must in this page
		totalNext := helper.CountTtlProductNxtPage(page, int(ctlgTotal))

		filter := bson.D{}

		if last_id != primitive.NilObjectID {
			id := bson.D{{Key: "$gt", Value: last_id}}
			filter = bson.D{{Key: "_id", Value: id}}
		}

		ids, err := findFirstAndLast(coll, filter, totalNext, opt)
		if err != nil {
			log.Print(rc.TaskIdentifier.Name, " Err first and last, ", err)
			return
		}

		idTable[i] = ids[0]
		last_id = ids[1]
	}
	// log.Println(idTable)
	err = cache.Pagination.StorePage(idTable)
	if err != nil {
		log.Println(err)
	}
}

// Deprecated
// try implement queryGetLastID, make it compatible with categories
func initLastID(last_id primitive.ObjectID, coll *mongo.Collection, opt *options.FindOptions) primitive.ObjectID {
	if last_id != primitive.NilObjectID {
		// skip 6 document
		skip := 6
		opt.SetLimit(int64(skip))
		filter := bson.D{{Key: "_id", Value: bson.D{{Key: "$gt", Value: last_id}}}}
		// Note
		// this function skip more than variable skip
		// because findFirstAndLast is skip again
		fl, err := findFirstAndLast(coll, filter, skip, opt)
		if err != nil {
			log.Println(err)
			return primitive.NilObjectID
		}

		return fl[1]
	}
	return primitive.NilObjectID
}

func queryGetLastID(last_id primitive.ObjectID, coll *mongo.Collection, additionalFilter bson.E, opt *options.FindOptions) primitive.ObjectID {
	if last_id == primitive.NilObjectID {
		return primitive.NilObjectID
	}

	// skip 6 document
	skip := 6
	opt.SetLimit(int64(skip))
	filter := bson.D{{Key: "_id", Value: bson.D{{Key: "$gt", Value: last_id}}}}
	if additionalFilter.Key != "" {
		filter = append(filter, additionalFilter)
	}
	// Note
	// this function skip more than variable skip
	// because findFirstAndLast is skip again
	fl, err := findFirstAndLast(coll, filter, skip, opt)
	if err != nil {
		log.Println(err)
		return primitive.NilObjectID
	}

	return fl[1]
}

func initLastIDCategory(last_id primitive.ObjectID, coll *mongo.Collection, category string, opt *options.FindOptions) primitive.ObjectID {
	return queryGetLastID(last_id, coll, bson.E{Key: repository.CategoryField, Value: category}, opt)
}

// using query $gt/greater than specific id
// directly skips the id that was thrown
//
// if in the database last line but item leng != content limit
// then query to the first line in database
//
// return id from first and last result
func findFirstAndLast(coll *mongo.Collection, filter bson.D, contentLimit int, opt *options.FindOptions) ([2]primitive.ObjectID, error) {
	var ids [2]primitive.ObjectID
	opt.SetLimit(int64(contentLimit))
	curr, err := coll.Find(context.TODO(), filter, opt)
	if err != nil {
		return ids, err
	}

	var decodeLeng int
	ids, decodeLeng, err = decodeFirstAndLastID(context.TODO(), curr)
	if err != nil {
		if contentLimit < 1 {
			return ids, err
		}
	}

	contentLimit -= decodeLeng
	log.Print("deviation: ", contentLimit)
	if contentLimit < 1 {
		log.Println(ids)
		return ids, nil
	}

	opt.SetLimit(int64(contentLimit))
	curr, err = coll.Find(context.TODO(), bson.D{}, opt)
	if err != nil {
		return ids, err
	}

	var ids2 [2]primitive.ObjectID
	ids2, _, err = decodeFirstAndLastID(context.TODO(), curr)
	if err != nil {
		return ids, err
	}

	if ids[0] == primitive.NilObjectID {
		ids[0] = ids2[1]
	}

	ids[1] = ids2[1]

	log.Println(ids)
	return ids, err
}

// [2]primitive.ObjectID first and las id
// int leng
func decodeFirstAndLastID(ctx context.Context, curr *mongo.Cursor) ([2]primitive.ObjectID, int, error) {
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
				return result, i, err
			}
			result[1] = tmp.Id
			break
		} else if i == 0 {
			err = curr.Decode(&tmp)
			if err != nil {
				return result, i, err
			}
			result[0] = tmp.Id
		}
		i++
	}

	return result, i, nil
}
