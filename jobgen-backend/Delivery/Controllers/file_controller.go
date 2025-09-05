package controllers

import (
	domain "jobgen-backend/Domain"
	"net/http"

	"github.com/gin-gonic/gin"
)

type FileController struct {
	fileUsecase domain.IFileUsecase
}

func NewFileController(fileUsecase domain.IFileUsecase) *FileController {
	return &FileController{
		fileUsecase: fileUsecase,
	}
}

// DeleteFile deletes a file by its ID
// @Summary Delete a file
// @Description Deletes a file owned by the authenticated user
// @Tags Files
// @Accept json
// @Produce json
// @Param id path string true "File ID"
// @Success 200 {object} StandardResponse "File deleted successfully"
// @Failure 401 {object} StandardResponse "Unauthorized"
// @Failure 404 {object} StandardResponse "File not found"
// @Failure 500 {object} StandardResponse "Internal server error"
// @Security BearerAuth
// @Router /files/:id [delete]
func (fc *FileController) DeleteFile(c *gin.Context) {
	dbID := c.Param("id")
	userID := c.GetString("user_id")
	if userID == "" {
		UnauthorizedResponse(c, "user id not found in context")
		return
	}
	err := fc.fileUsecase.Delete(c.Request.Context(), dbID, userID)
	if err != nil {
		switch err {
		case domain.ErrFileNotFound:
			NotFoundResponse(c, "file not found")
		case domain.ErrUnauthorized:
			UnauthorizedResponse(c, "you are not allowed to delete this file")
		default:
			InternalErrorResponse(c, "failed to delete file")
		}
		return
	}

	c.Status(http.StatusOK)
}

// UploadProfile handles profile picture uploads
// @Summary Upload profile picture
// @Description Uploads a profile picture for the authenticated user
// @Tags Files
// @Accept multipart/form-data
// @Produce json
// @Param file formData file true "Profile picture file"
// @Success 200 {object} StandardResponse "File uploaded successfully"
// @Failure 400 {object} StandardResponse "Validation error"
// @Failure 401 {object} StandardResponse "Unauthorized"
// @Failure 500 {object} StandardResponse "Internal server error"
// @Security BearerAuth
// @Router /files/upload/profile [post]
func (fc *FileController) UploadProfile(c *gin.Context) {
	fc.uploadFile(c, "profile-pictures")
}

// UploadDocument handles document uploads
// @Summary Upload a document
// @Description Uploads a document for the authenticated user
// @Tags Files
// @Accept multipart/form-data
// @Produce json
// @Param file formData file true "Document file"
// @Success 200 {object} StandardResponse "File uploaded successfully"
// @Failure 400 {object} StandardResponse "Validation error"
// @Failure 401 {object} StandardResponse "Unauthorized"
// @Failure 500 {object} StandardResponse "Internal server error"
// @Security BearerAuth
// @Router /files/upload/document [post]
func (fc *FileController) UploadDocument(c *gin.Context) {
	fc.uploadFile(c, "documents")
}

// Internal helper for uploading a file to a bucket
func (fc *FileController) uploadFile(c *gin.Context, bucket string) {
	userID := c.GetString("user_id")
	if userID == "" {
		UnauthorizedResponse(c, "user id not found in context")
		return
	}

	fileHeader, err := c.FormFile("file")
	if err != nil {
		ValidationErrorResponse(c, err)
		return
	}

	file, err := fileHeader.Open()
	if err != nil {
		InternalErrorResponse(c, "failed to open file")
		return
	}
	defer file.Close()

	data := domain.File{
		UserID:     userID,
		BucketName: bucket,
		FileName:   fileHeader.Filename,
		Size:       fileHeader.Size,
	}

	uploadedFile, err := fc.fileUsecase.Upload(c.Request.Context(), file, &data)
	if err != nil {
		switch err {
		case domain.ErrInvalidFileFormat, domain.ErrFileTooBig, domain.ErrUnknownFileType:
			ValidationErrorResponse(c, err)
		default:
			InternalErrorResponse(c, "failed to upload file")
		}
		return
	}

	SuccessResponse(c, http.StatusOK, "file uploaded successfully", uploadedFile)
}

// DownloadFile generates a presigned URL for downloading a file
// @Summary Download a file
// @Description Generates a presigned URL to download a file owned by the authenticated user
// @Tags Files
// @Accept json
// @Produce json
// @Param id path string true "File ID"
// @Success 200 {string} string "Presigned URL for the file"
// @Failure 401 {object} StandardResponse "Unauthorized"
// @Failure 404 {object} StandardResponse "File not found"
// @Failure 500 {object} StandardResponse "Internal server error"
// @Security BearerAuth
// @Router /files/:id [get]
func (fc *FileController) DownloadFile(c *gin.Context) {
	userID := c.GetString("user_id")
	if userID == "" {
		UnauthorizedResponse(c, "user id not found in context")
		return
	}
	dbID := c.Param("id")
	url, err := fc.fileUsecase.Download(c.Request.Context(), dbID, userID)
	if err != nil {
		switch err {
		case domain.ErrFileNotFound:
			NotFoundResponse(c, "file not found")
		case domain.ErrUnauthorized:
			UnauthorizedResponse(c, "you are not allowed to download this file")
		default:
			InternalErrorResponse(c, "failed to generate download link")
		}
		return
	}

	SuccessResponse(c, http.StatusOK, url, nil)
}

// GetMyProfilePicture returns a presigned URL for the current user's profile picture
// @Summary Get my profile picture
// @Description Fetch the current authenticated user's profile picture presigned URL
// @Tags Files
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} StandardResponse "Profile picture URL fetched successfully"
// @Failure 401 {object} StandardResponse "Unauthorized – user ID not found in context"
// @Failure 404 {object} StandardResponse "Profile picture not found"
// @Failure 500 {object} StandardResponse "Internal server error"
// @Router /files/profile-picture/me [get]
func (fc *FileController) GetMyProfilePicture(c *gin.Context) {
	userID := c.GetString("user_id")
	if userID == "" {
		UnauthorizedResponse(c, "user ID not found in context")
		return
	}

	url, err := fc.fileUsecase.GetProfilePictureByUserID(c.Request.Context(), userID)
	if err != nil {
		switch err {
		case domain.ErrFileNotFound:
			NotFoundResponse(c, "profile picture not found")
		default:
			InternalErrorResponse(c, "failed to fetch profile picture")
		}
		return
	}

	SuccessResponse(c, http.StatusOK, url, nil)
}

// GetProfilePicture returns a presigned URL for a user's profile picture
// @Summary Get profile picture
// @Description Fetch the profile picture presigned URL for a user if it exists
// @Tags Files
// @Accept json
// @Produce json
// @Param id path string true "User ID"
// @Success 200 {object} StandardResponse "Presigned URL returned successfully"
// @Failure 404 {object} StandardResponse "Profile picture not found"
// @Failure 500 {object} StandardResponse "Internal server error"
// @Router /files/profile-picture/:id [get]
func (fc *FileController) GetProfilePicture(c *gin.Context) {
	userID := c.Param("id")
	if userID == "" {
		UnauthorizedResponse(c, "user ID is required")
		return
	}

	url, err := fc.fileUsecase.GetProfilePictureByUserID(c.Request.Context(), userID)
	if err != nil {
		switch err {
		case domain.ErrFileNotFound:
			NotFoundResponse(c, "profile picture not found")
		default:
			InternalErrorResponse(c, "failed to fetch profile picture")
		}
		return
	}

	SuccessResponse(c, http.StatusOK, url, nil)
}
