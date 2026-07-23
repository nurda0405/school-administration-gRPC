package handlers

import (
	"context"
	"grpcapi/internals/models"
	"grpcapi/internals/repositories/mongodb"
	"grpcapi/pkg/utils"
	pb "grpcapi/proto/gen"
	"reflect"

	"go.mongodb.org/mongo-driver/v2/bson"
)

func (s *Server) AddTeachers(ctx context.Context, req *pb.Teachers) (*pb.Teachers, error) {
	client, err := mongodb.CreateMongoClient(ctx)
	if err != nil {
		return nil, utils.ErrorHandler(err, "Internal Error")
	}
	defer client.Disconnect(ctx)

	newTeachers := make([]*models.Teacher, len(req.GetTeachers()))
	for i, pbTeacher := range req.GetTeachers() {
		modelTeacher := models.Teacher{}
		pbVal := reflect.ValueOf(pbTeacher).Elem()
		modelVal := reflect.ValueOf(&modelTeacher).Elem()

		for i := 0; i < pbVal.NumField(); i++ {
			pbField := pbVal.Field(i)
			fieldName := pbVal.Type().Field(i).Name

			modelField := modelVal.FieldByName(fieldName)
			if modelField.IsValid() && modelField.CanSet() {
				modelField.Set(pbField)
			}
		}
		newTeachers[i] = &modelTeacher
	}

	var addedTeachers []*pb.Teacher
	for _, newTeacher := range newTeachers {
		res, err := client.Database("school").Collection("teachers").InsertOne(ctx, newTeacher)
		if err != nil {
			return nil, utils.ErrorHandler(err, "Error adding teacher to database")
		}
		objectID, ok := res.InsertedID.(bson.ObjectID)
		if ok {
			newTeacher.Id = objectID.Hex()
		}

		pbTeacher := &pb.Teacher{}
		modelVal := reflect.ValueOf(newTeacher).Elem()
		pbVal := reflect.ValueOf(pbTeacher).Elem()

		for i := 0; i < modelVal.NumField(); i++ {
			modelField := modelVal.Field(i)
			modelFieldType := modelVal.Type().Field(i)
			pbField := pbVal.FieldByName(modelFieldType.Name)

			if pbField.IsValid() && pbField.CanSet() {
				pbField.Set(modelField)
			}
		}
		addedTeachers = append(addedTeachers, pbTeacher)
	}

	return &pb.Teachers{Teachers: addedTeachers}, nil
}
