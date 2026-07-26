package handlers

import (
	"fmt"
	"grpcapi/pkg/utils"
	pb "grpcapi/proto/gen"
	"reflect"
	"strings"

	"go.mongodb.org/mongo-driver/v2/bson"
)

func buildFilter(object interface{}, model interface{}) (bson.M, error) {
	filter := bson.M{}

	if object == nil || reflect.ValueOf(object).IsNil() {
		return filter, nil
	}
	modelVal := reflect.ValueOf(model).Elem()
	modelType := modelVal.Type()

	reqVal := reflect.ValueOf(object).Elem()
	reqType := reqVal.Type()

	for i := 0; i < reqVal.NumField(); i++ {
		fieldVal := reqVal.Field(i)
		fieldName := reqType.Field(i).Name

		if fieldVal.IsValid() && !fieldVal.IsZero() {
			modelField := modelVal.FieldByName(fieldName)
			if modelField.IsValid() && modelField.CanSet() {
				modelField.Set(fieldVal)
			}
		}
	}

	for i := 0; i < modelVal.NumField(); i++ {
		fieldVal := modelVal.Field(i)
		fieldName := modelType.Field(i).Name

		if fieldVal.IsValid() && !fieldVal.IsZero() {
			bsonTag := modelType.Field(i).Tag.Get("bson")
			bsonTag = strings.TrimSuffix(bsonTag, ",omitempty")
			if bsonTag == "_id" {
				objId, err := bson.ObjectIDFromHex(reqVal.FieldByName(fieldName).Interface().(string))
				if err != nil {
					return nil, utils.ErrorHandler(err, "Invalid ID")
				}
				filter[bsonTag] = objId
			} else {
				filter[bsonTag] = fieldVal.Interface().(string)
			}
		}
	}
	fmt.Println(filter)
	return filter, nil
}

func buildSortOptions(sortFields []*pb.SortField) bson.D {
	var sortOptions bson.D

	for _, sortField := range sortFields {
		order := 1
		if sortField.GetOrder() == pb.Order_DESC {
			order = -1
		}
		sortOptions = append(sortOptions, bson.E{Key: sortField.Field, Value: order})
	}
	fmt.Println(sortOptions)
	return sortOptions
}
