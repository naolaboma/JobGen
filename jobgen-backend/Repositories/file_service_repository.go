package repositories

import (
	"context"
	"errors"
	domain "jobgen-backend/Domain"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type FileRepository struct {
	collection *mongo.Collection
}

func NewFileRepository(db *mongo.Database) domain.IFileRepository {
	return &FileRepository{
		collection: db.Collection("files"),
	}
}

// Delete implements domain.IFileRepository.
func (f *FileRepository) Delete(ctx context.Context, ID string) error {
	objID, err := primitive.ObjectIDFromHex(ID)
	if err != nil {
		return err
	}
	filter := bson.M{
		"_id": objID,
	}
	result, err := f.collection.DeleteOne(ctx, filter, nil)
	if err != nil {
		return err
	}
	if result.DeletedCount == 0 {
		return errors.New("file not found")
	}
	return nil
}

// FindByID implements domain.IFileRepository.
func (f *FileRepository) FindByID(ctx context.Context, ID string) (*domain.File, error) {
	objID, err := primitive.ObjectIDFromHex(ID)
	if err != nil {
		return nil, err
	}
	filter := bson.M{"_id": objID}
	result := f.collection.FindOne(ctx, filter)
	if result.Err() != nil {
		if errors.Is(result.Err(), mongo.ErrNoDocuments) {
			return nil, domain.ErrFileNotFound
		}
		return nil, result.Err()
	}

	var fileMeta domain.File
	if err := result.Decode(&fileMeta); err != nil {
		return nil, err
	}
	return &fileMeta, nil
}

// Store implements domain.IFileRepository.
func (f *FileRepository) Store(ctx context.Context, file *domain.File) error {
	result, err := f.collection.InsertOne(ctx, file)
	if err != nil {
		return err
	}
	id, ok := result.InsertedID.(primitive.ObjectID)
	if !ok {
		return errors.New("failed to assert database id")
	}
	file.ID = id.Hex()
	return nil
}

// FindByUserID fetches the profile picture of a specific user
func (r *FileRepository) FindByUserID(ctx context.Context, userID string) (*domain.File, error) {
	var file domain.File
	filter := bson.M{
		"user_id":     userID,
		"bucket_name": "profile-pictures", // ensure only profile pictures
	}
	err := r.collection.FindOne(ctx, filter).Decode(&file)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, domain.ErrFileNotFound
		}
		return nil, err
	}
	return &file, nil
}
