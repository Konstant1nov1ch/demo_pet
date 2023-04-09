package db

import (
	"app/internal"
	"app/internal/apperror"
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
func (d *db) Create(ctx context.Context, task internal.Task) (string, error) {
	d.logger.Debug("create task")
	result, err := d.collection.InsertOne(ctx, task)
	if err != nil {
		return "", fmt.Errorf("failed to create task due to error: %v", err)
	}
	d.logger.Debug("convert InsertedID to ObjectID")
	oid, ok := result.InsertedID.(primitive.ObjectID)
	if ok {
		return oid.Hex(), nil
	}
	d.logger.Trace(task)
	return "", fmt.Errorf("failed to convert objectid to hex. probably oid: %s", oid)
}
func (d *db) FindAll(ctx context.Context) (t []internal.Task, err error) {
	cursor, err := d.collection.Find(ctx, bson.M{})
	if cursor.Err() != nil {
		return t, fmt.Errorf("failed to find all tasks due to error: %v", err)
	}

	if err = cursor.All(ctx, &t); err != nil {
		return t, fmt.Errorf("failed to read all documents from cursor. error: %v", err)
	}

	return t, nil
}

func (d *db) FindOne(ctx context.Context, id string) (t internal.Task, err error) {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return t, fmt.Errorf("failed to convert hex to objectid. hex: %s", id)
	}

	filter := bson.M{"_id": oid}

	result := d.collection.FindOne(ctx, filter)
	if result.Err() != nil {
		if errors.Is(result.Err(), mongo.ErrNoDocuments) {
			return t, apperror.ErrNotFound
		}
		return t, fmt.Errorf("failed to find one task by id: %s due to error: %v", id, err)
	}
	if err = result.Decode(&t); err != nil {
		return t, fmt.Errorf("failed to decode task (id:%s) from DB due to error: %v", id, err)
	}
	return t, nil
}

func (d *db) UpdateTask(ctx context.Context, task internal.Task) error {
	objectID, err := primitive.ObjectIDFromHex(task.ID)
	if err != nil {
		return fmt.Errorf("failed to convert task ID to ObjectID. ID=%s", task.ID)
	}

	filter := bson.M{"_id": objectID}

	userBytes, err := bson.Marshal(task)
	if err != nil {
		return fmt.Errorf("failed to marhsal task. error: %v", err)
	}

	var updateUserObj bson.M
	err = bson.Unmarshal(userBytes, &updateUserObj)
	if err != nil {
		return fmt.Errorf("failed to unmarshal task bytes. error: %v", err)
	}

	delete(updateUserObj, "_id")

	update := bson.M{
		"$set": updateUserObj,
	}

	result, err := d.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return fmt.Errorf("failed to execute update task query. error: %v", err)
	}

	if result.MatchedCount == 0 {
		return apperror.ErrNotFound
	}

	d.logger.Tracef("Matched %d documents and Modified %d documents", result.MatchedCount, result.ModifiedCount)

	return nil

}

func (d *db) DeleteTask(ctx context.Context, id string) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return fmt.Errorf("failed to convert task ID to ObjectID. ID=%s", id)
	}

	filter := bson.M{"_id": objectID}

	result, err := d.collection.DeleteOne(ctx, filter)
	if err != nil {
		return fmt.Errorf("failed to execute query. error: %v", err)
	}
	if result.DeletedCount == 0 {
		return apperror.ErrNotFound
	}
	d.logger.Tracef("Deleted %d documents", result.DeletedCount)

	return nil
}

func NewStorage(database *mongo.Database, collection string, logger *logging.Logger) internal.Storage {
	return &db{
		collection: database.Collection(collection),
		logger:     logger,
	}
}
