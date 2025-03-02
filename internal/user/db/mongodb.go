package db

import (
	"context"
	"errors"
	"fmt"
	"notes-go/internal/user"
	"notes-go/pkg/logging"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type db struct {
	collection *mongo.Collection
	logger     *logging.Logger
}

// Create implements user.Storage.
func (d *db) Create(ctx context.Context, user user.User) (string, error) {
	d.logger.Debug("create user")
	result, err := d.collection.InsertOne(ctx, user)

	if err != nil {
		return "", fmt.Errorf("failed to create user due to error: %v", err)
	}

	d.logger.Debug("convert insertID to ObjectID")
	oid, ok := result.InsertedID.(primitive.ObjectID)
	
	if ok {
		return oid.Hex(), nil
	}

	d.logger.Trace(user)
	return "", fmt.Errorf("failed to convert objectid to hex")
}

// Delete implements user.Storage.
func (d *db) Delete(ctx context.Context, id string) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return fmt.Errorf("failed to convert user ID to ObjectId: %s", id)
	}

	filter := bson.M{"_id": objectID}

	result, err := d.collection.DeleteOne(ctx, filter)

	if err != nil {
		return fmt.Errorf("failed to execute query: %v", err)
	}

	if result.DeletedCount == 0 {
		return fmt.Errorf("not found")
	}

	d.logger.Tracef("Delete %d documents", result.DeletedCount)

	return nil
}

// FindOne implements user.Storage.
func (d *db) FindOne(ctx context.Context, id string) (u []user.User, err error) {
	oid, err := primitive.ObjectIDFromHex(id)

	if err != nil {
		return u, fmt.Errorf("failed to convert hex to objectId. hex: %s", id)
	}

	filter := bson.M{"_id": oid}
	result := d.collection.FindOne(ctx, filter)

	if result.Err() != nil {
		if errors.Is(result.Err(), mongo.ErrNoDocuments) {
			return u, fmt.Errorf("not found")
		}
		return u, fmt.Errorf("failed to find one user by id: %s due to error: %v", id, err)
	}

	if err = result.Decode(&u); err != nil {
		return u, fmt.Errorf("failed to decode user from db: %s - %v", id, err)
	}

	return u, nil
}

func (d *db) FindAll(ctx context.Context) (u user.User, err error) {
	cursor, err := d.collection.Find(ctx, bson.M{})

	if cursor.Err() != nil {
		return u, fmt.Errorf("failed to find all users: %v", err)
	}

	if err = cursor.All(ctx, &u); err != nil {
		return u, fmt.Errorf("failed to read all documents from curosr: %v", err)
	}

	return u, nil
}

// Update implements user.Storage.
func (d *db) Update(ctx context.Context, user user.User) error {
	objectID, err := primitive.ObjectIDFromHex(user.ID)
	if err != nil {
		return fmt.Errorf("failed to convert user ID to ObjectId: %s", user.ID)
	}

	filter := bson.M{"_id": objectID}
	userBytes, err := bson.Marshal(user)

	if err != nil {
		return fmt.Errorf("failed to marshal user: %v", err)
	}

	var updateUserObj bson.M
	err = bson.Unmarshal(userBytes, &updateUserObj)

	if err != nil {
		return fmt.Errorf("failed to unmarshal user bytes: %v", err)
	}

	delete(updateUserObj, "_id")
	update := bson.M{"$set": updateUserObj}

	result, err := d.collection.UpdateOne(ctx, filter, update)

	if err != nil {
		return fmt.Errorf("failed to execute update user query: %v", err)
	}

	if result.MatchedCount == 0 {
		return fmt.Errorf("not found")
	}

	d.logger.Tracef("Matched %d documents and modified %d documents", result.MatchedCount, result.ModifiedCount)
	return nil
}

func NewStorage(database *mongo.Database, collection string, logger *logging.Logger) user.Storage {
	return &db{
		collection: database.Collection(collection),
		logger:     logger,
	}
}
