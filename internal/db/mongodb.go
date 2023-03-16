package db

import (
	"app/internal"
	"app/pkg/logging"
	"context"
	"errors"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type db struct {
	collection *mongo.Collection
	logger     *logging.Logger
}

// nCtx, cancel := context.WithTimeout(ctx1, 1*time.Second)
func (d *db) CreateTask(ctx context.Context, task internal.ToDoList) (taskID string, err error) {
	d.logger.Debug("create task")
	result, err := d.collection.InsertOne(ctx, task)
	if err != nil {
		return "", fmt.Errorf("failed to create list: %v", err)
	}
	d.logger.Debug("convert insertedID to oid")
	oid, ok := result.InsertedID.(primitive.ObjectID)
	if ok {
		return oid.Hex(), nil
	}

	d.logger.Debug(task)
	return "", fmt.Errorf("failed to convert oid to hex: %s", oid)
}

func (d *db) FindOne(ctx context.Context, id string) (task internal.ToDoList, err error) {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return task, fmt.Errorf("faild  to convert hex to taskID:%s", id)
	}
	filter := bson.M{"_id": oid}
	result := d.collection.FindOne(ctx, filter)
	if result.Err() != nil {
		if errors.Is(result.Err(), mongo.ErrNoDocuments) {
			//TODO ErrEntityNotFound
			return task, fmt.Errorf("not found")
		}
		return task, fmt.Errorf("faild to find task by id: %s due to error: %v", id, err)
	}
	if err = result.Decode(&task); err != nil {
		return task, fmt.Errorf("faild to decode  task(id: %s) from DB due to error: %v", id, err)
	}
	return task, err
}

func (d *db) UpdateTask(ctx context.Context, task internal.ToDoList) error {
	objectID, objConvErr := primitive.ObjectIDFromHex(task.ID)
	if objConvErr != nil {
		return fmt.Errorf("faild to convert user ID to objectID. ID=%s", task.ID)
	}
	filter := bson.M{"id": objectID}

	taskBytes, taskMarshalErr := bson.Marshal(task)
	if taskMarshalErr != nil {
		return fmt.Errorf("faild to marshal task. error: %v", taskMarshalErr)
	}

	var updateTaskObj bson.M
	taskUnmarshalErr := bson.Unmarshal(taskBytes, &updateTaskObj)
	if taskUnmarshalErr != nil {
		return fmt.Errorf("faild to unmarshal task bytes. error: %v", taskUnmarshalErr)
	}
	delete(updateTaskObj, "_id")

	update := bson.M{
		"$set": updateTaskObj,
	}
	result, err := d.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return fmt.Errorf("faild to execute update task query. error:%v", err)
	}
	if result.MatchedCount == 0 {
		// TODO ErrEntityNotFound
		return fmt.Errorf("not found")
	}
	d.logger.Tracef("Matched %d documentd and Modified %d documents ", result.MatchedCount, result.ModifiedCount)
	return nil
}

func (d *db) DeleteTask(ctx context.Context, id string) error {
	objectID, objConvErr := primitive.ObjectIDFromHex(id)
	if objConvErr != nil {
		return fmt.Errorf("faild to convert user ID to objectID. ID=%s", id)
	}
	filter := bson.M{"id": objectID}

	result, err := d.collection.DeleteOne(ctx, filter)
	if err != nil {
		return fmt.Errorf("faild to execute quere")
	}
	if result.DeletedCount == 0 {
		//TODO ErrEntityNotFound
		return fmt.Errorf("not found")
	}
	d.logger.Tracef("Matched %d documentd ", result.DeletedCount)
	return nil
}

func NewStorage(database *mongo.Database, collection string, logger *logging.Logger) *db {
	return &db{
		collection: database.Collection(collection),
		logger:     logger,
	}
}
