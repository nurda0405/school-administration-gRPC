package mongodb

import (
	"context"
	"errors"
	"fmt"
	"grpcapi/internals/models"
	"grpcapi/pkg/utils"
	pb "grpcapi/proto/gen"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

func AddTeachersToDb(ctx context.Context, teachersFromReq []*pb.Teacher) ([]*pb.Teacher, error) {
	client, err := CreateMongoClient(ctx)
	if err != nil {
		return nil, utils.ErrorHandler(err, "Internal Error")
	}
	defer client.Disconnect(ctx)

	newTeachers := make([]*models.Teacher, len(teachersFromReq))
	for i, pbTeacher := range teachersFromReq {
		modelTeacher := mapPbToModel(pbTeacher, func() *models.Teacher { return &models.Teacher{} })
		newTeachers[i] = modelTeacher
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

		pbTeacher := mapModelToPb(newTeacher, func() *pb.Teacher { return &pb.Teacher{} })
		addedTeachers = append(addedTeachers, pbTeacher)
	}
	return addedTeachers, nil
}

func GetTeachersFromDb(ctx context.Context, sortOptions bson.D, filter bson.M) ([]*pb.Teacher, error) {
	client, err := CreateMongoClient(ctx)
	if err != nil {
		return nil, utils.ErrorHandler(err, "Internal Error")
	}
	defer client.Disconnect(ctx)

	coll := client.Database("school").Collection("teachers")
	var cursor *mongo.Cursor
	if len(sortOptions) > 0 {
		cursor, err = coll.Find(ctx, filter, options.Find().SetSort(sortOptions))
	} else {
		cursor, err = coll.Find(ctx, filter)
	}

	if err != nil {
		return nil, utils.ErrorHandler(err, "Internal Error")
	}
	defer cursor.Close(ctx)

	teachers, err := decodeEntities(ctx, cursor,
		func() *pb.Teacher {
			return &pb.Teacher{}
		},
		func() *models.Teacher {
			return &models.Teacher{}
		})
	if err != nil {
		return nil, utils.ErrorHandler(err, "Internal Error")
	}
	return teachers, nil
}

func ModifyTeachersInDb(ctx context.Context, pbTeachers []*pb.Teacher) ([]*pb.Teacher, error) {
	client, err := CreateMongoClient(ctx)
	if err != nil {
		return nil, utils.ErrorHandler(err, "Internal Error")
	}

	defer client.Disconnect(ctx)

	var updatedTeachers []*pb.Teacher
	for _, teacher := range pbTeachers {
		if teacher.Id == "" {
			return nil, utils.ErrorHandler(errors.New("blank id"), "blank id")
		}
		modelTeacher := mapPbToModel(teacher, func() *models.Teacher { return &models.Teacher{} })
		objID, err := bson.ObjectIDFromHex(teacher.Id)
		if err != nil {
			return nil, utils.ErrorHandler(err, "Invalid ID")
		}

		modelDoc, err := bson.Marshal(modelTeacher)
		if err != nil {
			return nil, utils.ErrorHandler(err, "Internal Error")
		}
		var updateDoc bson.M
		err = bson.Unmarshal(modelDoc, &updateDoc)
		if err != nil {
			return nil, utils.ErrorHandler(err, "Internal Error")
		}

		delete(updateDoc, "_id")

		_, err = client.Database("school").Collection("teachers").UpdateOne(ctx, bson.M{"_id": objID}, bson.M{"$set": updateDoc})
		if err != nil {
			return nil, utils.ErrorHandler(err, fmt.Sprintf("Error updating teacher", teacher.Id))
		}
		updatedTeacher := mapModelToPb(modelTeacher, func() *pb.Teacher { return &pb.Teacher{} })
		updatedTeachers = append(updatedTeachers, updatedTeacher)
	}
	return updatedTeachers, nil
}
