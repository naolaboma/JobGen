package domain

import "time"

type JobStatus string

const (
	StatusPending    JobStatus = "Pending"
	StatusProcessing JobStatus = "Processing"
	StatusCompleted  JobStatus = "Completed"
	StatusFailed     JobStatus = "Failed"
)

// CV is the core domain model for a curriculum vitae and its processing job.
type CV struct {
	ID              string       `json:"id" bson:"_id"`
	UserID          string       `json:"userId" bson:"userId"`
	FileStorageID   string       `json:"fileStorageId" bson:"fileStorageId"`
	FileName        string       `json:"fileName" bson:"fileName"`
	Status          JobStatus    `json:"status" bson:"status"`
	ProcessingError string       `json:"processingError,omitempty" bson:"processingError,omitempty"`
	RawText         string       `json:"rawText,omitempty" bson:"rawText,omitempty"`
	ProfileSummary  string       `json:"profileSummary,omitempty" bson:"profileSummary,omitempty"`
	Experiences     []Experience `json:"experiences,omitempty" bson:"experiences,omitempty"`
	Educations      []Education  `json:"educations,omitempty" bson:"educations,omitempty"`
	Skills          []string     `json:"skills,omitempty" bson:"skills,omitempty"`
	Suggestions     []Suggestion `json:"suggestions,omitempty" bson:"suggestions,omitempty"`
	Score           int          `json:"score" bson:"score"`
	CreatedAt       time.Time    `json:"createdAt" bson:"createdAt"`
	UpdatedAt       time.Time    `json:"updatedAt" bson:"updatedAt"`
}

type Experience struct {
	ID          string     `json:"id" bson:"id"`
	Title       string     `json:"title" bson:"title"`
	Company     string     `json:"company" bson:"company"`
	Location    string     `json:"location" bson:"location"`
	StartDate   time.Time  `json:"startDate" bson:"startDate"`
	EndDate     *time.Time `json:"endDate,omitempty" bson:"endDate,omitempty"`
	Description string     `json:"description" bson:"description"`
}

type Education struct {
	ID             string    `json:"id" bson:"id"`
	Degree         string    `json:"degree" bson:"degree"`
	Institution    string    `json:"institution" bson:"institution"`
	Location       string    `json:"location" bson:"location"`
	GraduationDate time.Time `json:"graduationDate" bson:"graduationDate"`
}

type Suggestion struct {
	ID      string `json:"id" bson:"id"`
	Type    string `json:"type" bson:"type"`
	Content string `json:"content" bson:"content"`
	Applied bool   `json:"applied" bson:"applied"`
}

// CVRepository defines the interface for CV data persistence.
type CVRepository interface {
	Create(cv *CV) error
	GetByID(id string) (*CV, error)
	UpdateStatus(id string, status JobStatus, procError ...string) error
	UpdateWithResults(id string, results *CV) error
}

// AIService defines the interface for the AI team's service.
type AIService interface {
	AnalyzeCV(rawText string) ([]Suggestion, error)
}
