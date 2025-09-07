package controllers

import (
	usecases "jobgen-backend/Usecases"
	"net/http"

	"github.com/gin-gonic/gin"
)

type CVController struct {
	cvUsecase usecases.CVUsecase
}

func NewCVController(uc usecases.CVUsecase) *CVController {
	return &CVController{cvUsecase: uc}
}

type CVFileRefRequest struct {
	FileID string `json:"fileId"`
}

// @Summary Start CV parsing job (multipart)
// @Description Upload a CV PDF via multipart and start a parsing job.
// @Tags CV
// @Accept mpfd
// @Produce json
// @Security BearerAuth
// @Param file formData file true "CV PDF file"
// @Success 202 {object} map[string]interface{} "Job accepted"
// @Failure 400 {object} controllers.StandardResponse
// @Failure 401 {object} controllers.StandardResponse
// @Failure 500 {object} controllers.StandardResponse
// @Router /cv/parse [post]
func (ctrl *CVController) StartParsingJobHandler(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "form field 'file' is required"})
		return
	}

	// userID is set by auth middleware as "user_id"
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
		return
	}

	jobID, err := ctrl.cvUsecase.CreateParsingJob(userID.(string), file)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create parsing job", "details": err.Error()})
		return
	}

	c.JSON(http.StatusAccepted, gin.H{
		"message": "CV parsing job accepted.",
		"jobId":   jobID,
	})
}

// StartParsingJobFromRef accepts either a multipart upload or a JSON body with {"fileId": "..."}
// @Summary Start CV parsing job (file or reference)
// @Description Start a parsing job by uploading a file (multipart) or providing an existing fileId in JSON body.
// @Tags CV
// @Accept mpfd
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param file formData file false "CV PDF file"
// @Param request body controllers.CVFileRefRequest false "Provide when using existing fileId"
// @Success 202 {object} map[string]interface{} "Job accepted"
// @Failure 400 {object} controllers.StandardResponse
// @Failure 401 {object} controllers.StandardResponse
// @Failure 500 {object} controllers.StandardResponse
// @Router /cv [post]
func (ctrl *CVController) StartParsingJobFromRef(c *gin.Context) {
	// userID is set by auth middleware as "user_id"
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
		return
	}

	// Try multipart first
	if file, err := c.FormFile("file"); err == nil && file != nil {
		jobID, err := ctrl.cvUsecase.CreateParsingJob(userID.(string), file)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create parsing job", "details": err.Error()})
			return
		}
		c.JSON(http.StatusAccepted, gin.H{"message": "CV parsing job accepted.", "jobId": jobID})
		return
	}

	// Otherwise expect JSON with fileId
	var body struct {
		FileID string `json:"fileId" binding:"required"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "file or fileId is required"})
		return
	}
	jobID, err := ctrl.cvUsecase.CreateParsingJobFromFileID(userID.(string), body.FileID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create parsing job", "details": err.Error()})
		return
	}
	c.JSON(http.StatusAccepted, gin.H{"message": "CV parsing job accepted.", "jobId": jobID})
}

// GetParsingJobStatusHandler retrieves the status and result of a CV parsing job
// @Summary Get CV parsing job status
// @Description Fetch the current status and (if finished) the parsed result of a CV parsing job.
// @Tags CV
// @Produce json
// @Security BearerAuth
// @Param jobId path string true "Parsing Job ID"
// @Success 200 {object} map[string]interface{} "Job status and result"
// @Failure 401 {object} controllers.StandardResponse "User not authenticated"
// @Failure 403 {object} controllers.StandardResponse "User not authorized to view this job"
// @Failure 404 {object} controllers.StandardResponse "Job not found"
// @Router /cv/{jobId} [get]
func (ctrl *CVController) GetParsingJobStatusHandler(c *gin.Context) {
	jobID := c.Param("jobId")
	cv, err := ctrl.cvUsecase.GetJobStatusAndResult(jobID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "job not found"})
		return
	}

	// Basic authorization check: does the requesting user own this CV?
	requestingUserID, _ := c.Get("user_id")
	if cv.UserID != requestingUserID {
		c.JSON(http.StatusForbidden, gin.H{"error": "you are not authorized to view this CV"})
		return
	}

	c.JSON(http.StatusOK, cv)
}
