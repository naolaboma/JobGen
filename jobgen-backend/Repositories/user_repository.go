package repositories

import (
	"context"
	"fmt"
	domain "jobgen-backend/Domain"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type UserRepository struct {
	collection *mongo.Collection
}

func NewUserRepository(db *mongo.Database) domain.IUserRepository {
	repo := &UserRepository{
		collection: db.Collection("users"),
	}
	
	// Create indexes for email and username uniqueness
	repo.createIndexes()
	
	return repo
}

func (r *UserRepository) createIndexes() {
	ctx := context.Background()
	
	// Create unique index for email
	emailIndex := mongo.IndexModel{
		Keys:    bson.D{{Key: "email", Value: 1}},
		Options: options.Index().SetUnique(true),
	}
	
	// Create unique index for username
	usernameIndex := mongo.IndexModel{
		Keys:    bson.D{{Key: "username", Value: 1}},
		Options: options.Index().SetUnique(true),
	}
	
	r.collection.Indexes().CreateMany(ctx, []mongo.IndexModel{emailIndex, usernameIndex})
}

func (r *UserRepository) Create(ctx context.Context, user *domain.User) error {
	// Check if email already exists
	existingUser, err := r.GetByEmail(ctx, user.Email)
	if err != nil && err != domain.ErrUserNotFound {
		return err
	}
	if existingUser != nil {
		return domain.ErrEmailTaken
	}

	// Check if username already exists
	existingUser, err = r.GetByUsername(ctx, user.Username)
	if err != nil && err != domain.ErrUserNotFound {
		return err
	}
	if existingUser != nil {
		return domain.ErrUsernameTaken
	}

	// Generate new ObjectID
	user.ID = primitive.NewObjectID().Hex()
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()
	
	_, err = r.collection.InsertOne(ctx, user)
	if err != nil {
		// Handle MongoDB duplicate key errors
		if mongo.IsDuplicateKeyError(err) {
			if err.Error() == "email" {
				return domain.ErrEmailTaken
			}
			return domain.ErrUsernameTaken
		}
		return err
	}
	
	return nil
}

func (r *UserRepository) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
    var user domain.User
    email = strings.ToLower(strings.TrimSpace(email))

    filter := bson.M{"email": email}
    fmt.Println("Looking for user with filter:", filter, "in collection:", r.collection.Name())

    res := r.collection.FindOne(ctx, filter)
    if res.Err() != nil {
        if res.Err() == mongo.ErrNoDocuments {
            return nil, domain.ErrUserNotFound
        }
        return nil, fmt.Errorf("query failed: %w", res.Err())
    }

    if err := res.Decode(&user); err != nil {
        return nil, fmt.Errorf("failed to decode user: %w", err)
    }

    return &user, nil
}

func (r *UserRepository) GetByID(ctx context.Context, id string) (*domain.User, error) {
    var user domain.User
    objectID, err := primitive.ObjectIDFromHex(id)
    if err != nil {
        return nil, domain.ErrUserNotFound
    }
    filter := bson.M{"_id": objectID}
    fmt.Println("User that is needed to be found", filter)
    err = r.collection.FindOne(ctx, filter).Decode(&user)
    if err == mongo.ErrNoDocuments {
        return nil, domain.ErrUserNotFound
    }
    return &user, err
}

func (r *UserRepository) GetByUsername(ctx context.Context, username string) (*domain.User, error) {
	var user domain.User
	filter := bson.M{"username": username}
	err := r.collection.FindOne(ctx, filter).Decode(&user)
	if err == mongo.ErrNoDocuments {
		return nil, domain.ErrUserNotFound
	}
	return &user, err
}

func (r *UserRepository) Update(ctx context.Context, user *domain.User) error {
    user.UpdatedAt = time.Now()
    filter := bson.M{"_id": user.ID} // use string, not ObjectID
    update := bson.M{"$set": user}

    result, err := r.collection.UpdateOne(ctx, filter, update)
    if err != nil {
        return err
    }

    if result.MatchedCount == 0 {
        return domain.ErrUserNotFound
    }

    return nil
}


func (r *UserRepository) UpdatePassword(ctx context.Context, userID, newPasswordHash string) error {
    objectID, err := primitive.ObjectIDFromHex(userID)
    if err != nil {
        return domain.ErrUserNotFound
    }
    filter := bson.M{"_id": objectID}
    update := bson.M{
        "$set": bson.M{
            "password":   newPasswordHash,
            "updated_at": time.Now(),
        },
    }
    
    result, err := r.collection.UpdateOne(ctx, filter, update)
    if err != nil {
        return err
    }
    
    if result.MatchedCount == 0 {
        return domain.ErrUserNotFound
    }
    
    return nil
}

func (r *UserRepository) UpdateLastLogin(ctx context.Context, userID string) error {
    objectID, err := primitive.ObjectIDFromHex(userID)
    if err != nil {
        return domain.ErrUserNotFound
    }
    filter := bson.M{"_id": objectID}
    now := time.Now()
    update := bson.M{"$set": bson.M{"last_login_at": now}}
    
    _, err = r.collection.UpdateOne(ctx, filter, update)
    return err
}

func (r *UserRepository) Delete(ctx context.Context, userID string) error {
    objectID, err := primitive.ObjectIDFromHex(userID)
    if err != nil {
        return domain.ErrUserNotFound
    }
    filter := bson.M{"_id": objectID}
    result, err := r.collection.DeleteOne(ctx, filter)
    if err != nil {
        return err
    }
    
    if result.DeletedCount == 0 {
        return domain.ErrUserNotFound
    }
    
    return nil
}

func (r *UserRepository) List(ctx context.Context, filter domain.UserFilter) ([]domain.User, int64, error) {
	// Build MongoDB filter
	mongoFilter := bson.M{}
	
	if filter.Role != nil {
		mongoFilter["role"] = *filter.Role
	}
	
	if filter.IsActive != nil {
		mongoFilter["is_active"] = *filter.IsActive
	}
	
	if filter.Search != "" {
		// Search in email, username, or full_name
		mongoFilter["$or"] = []bson.M{
			{"email": bson.M{"$regex": filter.Search, "$options": "i"}},
			{"username": bson.M{"$regex": filter.Search, "$options": "i"}},
			{"full_name": bson.M{"$regex": filter.Search, "$options": "i"}},
		}
	}
	
	// Count total documents
	total, err := r.collection.CountDocuments(ctx, mongoFilter)
	if err != nil {
		return nil, 0, err
	}
	
	// Build sort options
	sortOrder := 1
	if filter.SortOrder == "desc" {
		sortOrder = -1
	}
	
	sortBy := "created_at"
	if filter.SortBy != "" {
		sortBy = filter.SortBy
	}
	
	// Calculate pagination
	skip := (filter.Page - 1) * filter.Limit
	
	// Find options
	findOptions := options.Find().
		SetSort(bson.D{{Key: sortBy, Value: sortOrder}}).
		SetSkip(int64(skip)).
		SetLimit(int64(filter.Limit))
	
	cursor, err := r.collection.Find(ctx, mongoFilter, findOptions)
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(ctx)
	
	var users []domain.User
	if err := cursor.All(ctx, &users); err != nil {
		return nil, 0, err
	}
	
	return users, total, nil
}

func (r *UserRepository) UpdateRole(ctx context.Context, userID string, role domain.Role) error {
    objectID, err := primitive.ObjectIDFromHex(userID)
    if err != nil {
        return domain.ErrUserNotFound
    }
    filter := bson.M{"_id": objectID}
    update := bson.M{
        "$set": bson.M{
            "role":       role,
            "updated_at": time.Now(),
        },
    }
    
    result, err := r.collection.UpdateOne(ctx, filter, update)
    if err != nil {
        return err
	}
    
    if result.MatchedCount == 0 {
        return domain.ErrUserNotFound
    }
    
    return nil
}

func (r *UserRepository) ToggleActiveStatus(ctx context.Context, userID string) error {
    objectID, err := primitive.ObjectIDFromHex(userID)
    if err != nil {
        return domain.ErrUserNotFound
    }
    // First get the current user to know current status
    user, err := r.GetByID(ctx, userID)
    if err != nil {
        return err
    }
    
    // Toggle the status
    newStatus := !user.IsActive
    
    filter := bson.M{"_id": objectID}
    update := bson.M{
        "$set": bson.M{
            "is_active":  newStatus,
            "updated_at": time.Now(),
        },
    }
    
    result, err := r.collection.UpdateOne(ctx, filter, update)
    if err != nil {
        return err
    }
    
    if result.MatchedCount == 0 {
        return domain.ErrUserNotFound
    }
    
    return nil
}
