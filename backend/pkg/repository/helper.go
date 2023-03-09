package repository

import (
	"context"
	"crypto/rand"
	"fmt"
	"math/big"
	"reflect"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// Extract value to mongo document
// param must be struct with bson tag
// bson tag as key for element bson.E
func DocumentFromModel(model interface{}) bson.D {
	var doc bson.D
	cpVal := reflect.ValueOf(model)
	cpType := cpVal.Type()

	for i := 0; i < cpType.NumField(); i++ {
		// skip if insert tag equal 0
		if cpType.Field(i).Tag.Get("insert") == "0" {
			continue
		}

		structField := cpVal.Field(i)
		bsonKey := cpType.Field(i).Tag.Get("bson")

		doc = append(doc, bson.E{Key: bsonKey, Value: structField.Interface()})
	}

	return doc
}

func DecodeCatalogCursor(ctx context.Context, curr *mongo.Cursor) ([]DocCatalog, error) {
	var result []DocCatalog
	for curr.Next(ctx) {
		var tmp DocCatalog
		if err := curr.Decode(&tmp); err != nil {
			return result, err
		}
		result = append(result, tmp)
	}

	if len(result) < 1 {
		return result, mongo.ErrEmptySlice
	}
	return result, nil
}

func DecodeId(ctx context.Context, curr *mongo.Cursor) ([]primitive.ObjectID, error) {
	type IDDecoder struct {
		Id primitive.ObjectID `bson:"_id"`
	}

	var result []primitive.ObjectID
	for curr.Next(ctx) {
		var tmp IDDecoder
		if err := curr.Decode(&tmp); err != nil {
			return result, err
		}
		result = append(result, tmp.Id)
	}

	if len(result) < 1 {
		return result, mongo.ErrEmptySlice
	}

	return result, nil
}

func RandInt(min, max int32) int32 {
	i, _ := rand.Int(rand.Reader, big.NewInt(int64(max-min)))
	return min + int32(i.Int64())
}

// give random string with certain length
// random string generating alphabet
// min represent a, max represent z
func RandString(leng int) string {
	var preStr []rune
	var min, max int32 = 97, 122

	for i := 0; i < leng; i++ {
		generated_num := RandInt(min, max)
		preStr = append(preStr, generated_num)
	}

	return string(preStr)
}

// Create very random string
// using RandString to generate block string
func RandName(single ...bool) string {
	firstName := RandString(int(RandInt(4, 10)))
	for _, v := range single {
		if v {
			return firstName
		}
	}

	lastName := RandString(int(RandInt(4, 10)))
	return fmt.Sprintf("%s %s", firstName, lastName)
}
