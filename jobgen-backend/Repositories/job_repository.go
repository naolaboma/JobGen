package repositories

import (
	"context"
	"fmt"
	domain "jobgen-backend/Domain"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type JobRepository struct {
	collection *mongo.Collection
}

func NewJobRepository(db *mongo.Database) domain.IJobRepository {
	repo := &JobRepository{
		collection: db.Collection("jobs"),
	}
	
	// Create indexes for efficient querying
	repo.createIndexes()
	
	return repo
}

func (r *JobRepository) createIndexes() {
	ctx := context.Background()
	
	// Text index for search functionality
	textIndex := mongo.IndexModel{
		Keys: bson.D{
			{Key: "title", Value: "text"},
			{Key: "company_name", Value: "text"},
			{Key: "description", Value: "text"},
		},
	}
	
	// Unique index on apply_url to prevent duplicates
	applyURLIndex := mongo.IndexModel{
		Keys:    bson.D{{Key: "apply_url", Value: 1}},
		Options: options.Index().SetUnique(true),
	}
	
	// Index on extracted_skills for filtering
	skillsIndex := mongo.IndexModel{
		Keys: bson.D{{Key: "extracted_skills", Value: 1}},
	}
	
	// Index on posted_at for sorting
	postedAtIndex := mongo.IndexModel{
		Keys: bson.D{{Key: "posted_at", Value: -1}},
	}
	
	// Index on source for filtering
	sourceIndex := mongo.IndexModel{
		Keys: bson.D{{Key: "source", Value: 1}},
	}
	
	// Index on location for filtering
	locationIndex := mongo.IndexModel{
		Keys: bson.D{{Key: "location", Value: 1}},
	}
	
	// Index on created_at for sorting
	createdAtIndex := mongo.IndexModel{
		Keys: bson.D{{Key: "created_at", Value: -1}},
	}
	
	// Index on remote_ok_id for RemoteOK jobs
	remoteOKIndex := mongo.IndexModel{
		Keys: bson.D{{Key: "remote_ok_id", Value: 1}},
		Options: options.Index().SetSparse(true),
	}
	
	r.collection.Indexes().CreateMany(ctx, []mongo.IndexModel{
		textIndex,
		applyURLIndex,
		skillsIndex,
		postedAtIndex,
		sourceIndex,
		locationIndex,
		createdAtIndex,
		remoteOKIndex,
	})
}

func (r *JobRepository) Create(ctx context.Context, job *domain.Job) error {
	// Generate new ObjectID if not set
	if job.ID == "" {
		job.ID = primitive.NewObjectID().Hex()
	}
	
	now := time.Now()
	job.CreatedAt = now
	job.UpdatedAt = now
	
	_, err := r.collection.InsertOne(ctx, job)
	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			return fmt.Errorf("job with this apply URL already exists")
		}
		return err
	}
	
	return nil
}

func (r *JobRepository) GetByID(ctx context.Context, id string) (*domain.Job, error) {
    var job domain.Job
    
    // First try with string ID
    filter := bson.M{"_id": id}
    err := r.collection.FindOne(ctx, filter).Decode(&job)
    if err == nil {
        return &job, nil
    }
    
    // If not found or error, try with ObjectID
    objectID, err := primitive.ObjectIDFromHex(id)
    if err != nil {
        return nil, domain.ErrNotFound
    }
    
    filter = bson.M{"_id": objectID}
    err = r.collection.FindOne(ctx, filter).Decode(&job)
    if err == mongo.ErrNoDocuments {
        return nil, domain.ErrNotFound
    }
    
    return &job, err
}

func (r *JobRepository) GetByApplyURL(ctx context.Context, applyURL string) (*domain.Job, error) {
	var job domain.Job
	filter := bson.M{"apply_url": applyURL}
	err := r.collection.FindOne(ctx, filter).Decode(&job)
	if err == mongo.ErrNoDocuments {
		return nil, nil // Not found, but not an error for checking existence
	}
	return &job, err
}

func (r *JobRepository) Update(ctx context.Context, job *domain.Job) error {
    job.UpdatedAt = time.Now()
    
    // First try with string ID
    filter := bson.M{"_id": job.ID}
    update := bson.M{"$set": job}
    
    result, err := r.collection.UpdateOne(ctx, filter, update)
    if err == nil && result.MatchedCount > 0 {
        return nil
    }
    
    // If not found or error, try with ObjectID
    objectID, err := primitive.ObjectIDFromHex(job.ID)
    if err != nil {
        return domain.ErrNotFound
    }
    
    filter = bson.M{"_id": objectID}
    result, err = r.collection.UpdateOne(ctx, filter, update)
    if err != nil {
        return err
    }
    
    if result.MatchedCount == 0 {
        return domain.ErrNotFound
    }
    
    return nil
}

func (r *JobRepository) Delete(ctx context.Context, id string) error {
    // First try with string ID
    filter := bson.M{"_id": id}
    result, err := r.collection.DeleteOne(ctx, filter)
    if err == nil && result.DeletedCount > 0 {
        return nil
    }
    
    // If not found or error, try with ObjectID
    objectID, err := primitive.ObjectIDFromHex(id)
    if err != nil {
        return domain.ErrNotFound
    }
    
    filter = bson.M{"_id": objectID}
    result, err = r.collection.DeleteOne(ctx, filter)
    if err != nil {
        return err
    }
    
    if result.DeletedCount == 0 {
        return domain.ErrNotFound
    }
    
    return nil
}

func (r *JobRepository) List(ctx context.Context, filter domain.JobFilter) ([]domain.Job, int64, error) {
	// Build MongoDB filter
	mongoFilter := bson.M{}
	
	// Text search
	if filter.Query != "" {
		mongoFilter["$text"] = bson.M{"$search": filter.Query}
	}
	
	// Skills filter
	if len(filter.Skills) > 0 {
		mongoFilter["extracted_skills"] = bson.M{"$in": filter.Skills}
	}
	
	// Location filter
	if filter.Location != "" {
		mongoFilter["location"] = bson.M{"$regex": filter.Location, "$options": "i"}
	}
	
	// Sponsorship filter
	if filter.Sponsorship != nil {
		mongoFilter["is_sponsorship_available"] = *filter.Sponsorship
	}
	
	// Source filter
	if filter.Source != "" {
		mongoFilter["source"] = filter.Source
	}
	
	// Count total documents
	total, err := r.collection.CountDocuments(ctx, mongoFilter)
	if err != nil {
		return nil, 0, err
	}
	
	// Build sort options
	sortOrder := -1 // desc by default
	if filter.SortOrder == "asc" {
		sortOrder = 1
	}
	
	sortBy := "posted_at" // default sort by posted_at
	if filter.SortBy != "" {
		sortBy = filter.SortBy
	}
	
	// Calculate pagination
	skip := (filter.Page - 1) * filter.Limit
	
	// Build find options
	findOptions := options.Find().
		SetSort(bson.D{{Key: sortBy, Value: sortOrder}}).
		SetSkip(int64(skip)).
		SetLimit(int64(filter.Limit))
	
	// Add text score projection if text search is used
	if filter.Query != "" {
		findOptions.SetProjection(bson.M{"score": bson.M{"$meta": "textScore"}})
		findOptions.SetSort(bson.D{{Key: "score", Value: bson.M{"$meta": "textScore"}}})
	}
	
	cursor, err := r.collection.Find(ctx, mongoFilter, findOptions)
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(ctx)
	
	var jobs []domain.Job
	if err := cursor.All(ctx, &jobs); err != nil {
		return nil, 0, err
	}
	
	return jobs, total, nil
}

// todo try to review bulkUpsert again there is some problem inside creation
func (r *JobRepository) BulkUpsert(ctx context.Context, jobs []domain.Job) error {
	if len(jobs) == 0 {
		return nil
	}

	var operations []mongo.WriteModel
	now := time.Now()

	for _, job := range jobs {
		// Generate ID if not set
		if job.ID == "" {
			job.ID = primitive.NewObjectID().Hex()
		}

		// Ensure timestamps
		if job.CreatedAt.IsZero() {
			job.CreatedAt = now
		}
		job.UpdatedAt = now

		// Updatable fields (exclude _id, created_at)
		updateFields := bson.M{
			"title":                  job.Title,
			"company_name":           job.CompanyName,
			"location":               job.Location,
			"description":            job.Description,
			"full_description_html":  job.FullDescriptionHTML,
			"apply_url":              job.ApplyURL,
			"source":                 job.Source,
			"posted_at":              job.PostedAt,
			"is_sponsorship_available": job.IsSponsorshipAvailable,
			"extracted_skills":       job.ExtractedSkills,
			"remote_ok_id":           job.RemoteOKID,
			"salary":                 job.Salary,
			"tags":                   job.Tags,
			"company_logo":           job.CompanyLogo,
			"original_data":          job.OriginalData,
			"updated_at":             job.UpdatedAt,
		}

		update := bson.M{
			"$set": updateFields,
			"$setOnInsert": bson.M{
				"_id":        job.ID,
				"created_at": job.CreatedAt,
			},
		}

		// Filter by unique apply_url
		filter := bson.M{"apply_url": job.ApplyURL}

		operation := mongo.NewUpdateOneModel().
			SetFilter(filter).
			SetUpdate(update).
			SetUpsert(true)

		operations = append(operations, operation)
	}

	opts := options.BulkWrite().SetOrdered(false)
	_, err := r.collection.BulkWrite(ctx, operations, opts)
	if err != nil {
		return fmt.Errorf("bulk upsert failed: %w", err)
	}

	return nil
}

func (r *JobRepository) GetJobsForMatching(ctx context.Context, limit int, offset int) ([]domain.Job, error) {
	filter := bson.M{}
	
	findOptions := options.Find().
		SetSort(bson.D{{Key: "posted_at", Value: -1}}).
		SetSkip(int64(offset)).
		SetLimit(int64(limit))
	
	cursor, err := r.collection.Find(ctx, filter, findOptions)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)
	
	var jobs []domain.Job
	if err := cursor.All(ctx, &jobs); err != nil {
		return nil, err
	}
	
	return jobs, nil
}
