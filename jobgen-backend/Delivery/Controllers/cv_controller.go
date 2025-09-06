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

func (ctrl *CVController) StartParsingJobHandler(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "form field 'file' is required"})
		return
	}

	// Extract user ID set by AuthMiddleware
	userID, exists := c.Get("user_id")
	if !exists || userID.(string) == "" {
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

func (ctrl *CVController) GetParsingJobStatusHandler(c *gin.Context) {
	jobID := c.Param("jobId")
	cv, err := ctrl.cvUsecase.GetJobStatusAndResult(jobID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "job not found"})
		return
	}

	// Basic authorization check: does the requesting user own this CV?
	requestingUserID, _ := c.Get("user_id")
	if requestingUserID == nil || cv.UserID != requestingUserID.(string) {
		c.JSON(http.StatusForbidden, gin.H{"error": "you are not authorized to view this CV"})
		return
	}

	c.JSON(http.StatusOK, cv)
}
