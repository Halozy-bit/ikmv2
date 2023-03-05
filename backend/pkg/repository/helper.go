package repository

import (
	"context"
	"reflect"

	"go.mongodb.org/mongo-driver/bson"
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
