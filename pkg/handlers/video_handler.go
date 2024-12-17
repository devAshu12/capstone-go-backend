package handlers

import (
	"errors"
	"github/devAshu12/learning_platform_GO_backend/internal/utils"
	"github/devAshu12/learning_platform_GO_backend/pkg/config"
	"github/devAshu12/learning_platform_GO_backend/pkg/db"
	"github/devAshu12/learning_platform_GO_backend/pkg/models"
	"github/devAshu12/learning_platform_GO_backend/pkg/types"
	"net/http"

	"gorm.io/gorm"
)

func CreateVideo(w http.ResponseWriter, r *http.Request) {

	// Parse the multipart form with a max size limit (e.g., 10MB)
	err := r.ParseMultipartForm(10 << 20) // 10MB
	if err != nil {
		appErr := types.NewAppError(http.StatusBadRequest, "Unable to parse form", err)
		utils.RespondWithError(w, appErr)
		return
	}

	// Extract title and module_id
	title := r.FormValue("title")
	moduleID := r.FormValue("module_id")
	// Validate fields
	if title == "" || moduleID == "" {
		appErr := types.NewAppError(http.StatusBadRequest, "Missing required fields: title or module_id", errors.New("missing required fields: title or module_id"))
		utils.RespondWithError(w, appErr)
		return
	}
	// Extract file
	file, fileHeader, err := r.FormFile("file")
	if err != nil {
		appErr := types.NewAppError(http.StatusBadRequest, "Error retrieving the file", err)
		utils.RespondWithError(w, appErr)
		return
	}

	defer file.Close()

	if fileHeader.Size == 0 {
		appErr := types.NewAppError(http.StatusBadRequest, "Empty File", err)
		utils.RespondWithError(w, appErr)
		return
	}

	// Optional: Set a folder name for Cloudinary
	folder := "capstone-folder"

	// Upload the file to Cloudinary
	securedURL, publicID, err := config.UploadFile(file, fileHeader.Filename, folder)
	if err != nil {
		appErr := types.NewAppError(http.StatusInternalServerError, "Failed to upload file", err)
		utils.RespondWithError(w, appErr)
		return
	}

	video := models.Video{
		Title:    title,
		ModuleID: moduleID,
		PublicID: publicID,
		URL:      securedURL,
	}

	if err := db.DB.Create(&video).Error; err != nil {
		appErr := types.NewAppError(http.StatusInternalServerError, "Failed to upload video", err)
		utils.RespondWithError(w, appErr)
		return
	}

	types.RespondWithJSON(w, http.StatusCreated, video)
}

func RemoveVideo(w http.ResponseWriter, r *http.Request) {
	// Get video_id from query parameters
	videoID := r.URL.Query().Get("video_id")
	if videoID == "" {
		appErr := types.NewAppError(http.StatusBadRequest, "video_id is required", nil)
		utils.RespondWithError(w, appErr)
		return
	}

	// Initialize a variable to hold the video record
	var video models.Video

	// Search for the video by its ID in the database
	result := db.DB.Where("id = ?", videoID).First(&video)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			appErr := types.NewAppError(http.StatusNotFound, "Video not found", nil)
			utils.RespondWithError(w, appErr)
		} else {
			appErr := types.NewAppError(http.StatusInternalServerError, "Error fetching video", result.Error)
			utils.RespondWithError(w, appErr)
		}
		return
	}

	// Delete the file from Cloudinary using its public ID
	err := config.DeleteFile(video.PublicID)
	if err != nil {
		appErr := types.NewAppError(http.StatusInternalServerError, "Error deleting video from Cloudinary", err)
		utils.RespondWithError(w, appErr)
		return
	}

	// Delete the video record from the database
	result = db.DB.Delete(&video)
	if result.Error != nil {
		appErr := types.NewAppError(http.StatusInternalServerError, "Error deleting video from database", result.Error)
		utils.RespondWithError(w, appErr)
		return
	}

	// Respond with a success message
	types.RespondWithJSON(w, http.StatusOK, map[string]string{
		"message": "Video successfully deleted",
	})
}
