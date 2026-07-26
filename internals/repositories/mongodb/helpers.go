package mongodb

import (
	"context"
	"grpcapi/pkg/utils"
	"reflect"

	"go.mongodb.org/mongo-driver/v2/mongo"
)

func decodeEntities[T any, M any](ctx context.Context, cursor *mongo.Cursor, newEntity func() *T, newModel func() *M) ([]*T, error) {
	var entities []*T
	for cursor.Next(ctx) {
		model := newModel()
		var err = cursor.Decode(&model)
		if err != nil {
			return nil, utils.ErrorHandler(err, "Internal Error")
		}
		entity := newEntity()
		modelVal := reflect.ValueOf(model).Elem()
		pbVal := reflect.ValueOf(entity).Elem()

		for i := 0; i < modelVal.NumField(); i++ {
			modelField := modelVal.Field(i)
			modelFieldName := modelVal.Type().Field(i).Name

			pbField := pbVal.FieldByName(modelFieldName)
			if pbField.IsValid() && pbField.CanSet() {
				pbField.Set(modelField)
			}
		}

		entities = append(entities, entity)

	}

	err := cursor.Err()
	if err != nil {
		return nil, utils.ErrorHandler(err, "Internal Error")
	}
	return entities, nil
}
