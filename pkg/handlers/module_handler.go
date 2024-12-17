package handlers

import (
	"encoding/json"
	"fmt"
	"github/devAshu12/learning_platform_GO_backend/internal/middlewares"
	"github/devAshu12/learning_platform_GO_backend/internal/utils"
	"github/devAshu12/learning_platform_GO_backend/pkg/config"
	"github/devAshu12/learning_platform_GO_backend/pkg/db"
	"github/devAshu12/learning_platform_GO_backend/pkg/models"
	"github/devAshu12/learning_platform_GO_backend/pkg/types"
	"net/http"
)

func CreateModule(w http.ResponseWriter, r *http.Request) {
	var createModuleReq types.CreateModuleReq
	err := json.NewDecoder(r.Body).Decode(&createModuleReq)
	if err != nil {
		appErr := types.NewAppError(http.StatusBadRequest, "Invalid request format", err)
		utils.RespondWithError(w, appErr)
		return
	}

	// Validate the input request
	if err := config.ValidateRequest(createModuleReq); err != nil {
		appErr := types.NewAppError(http.StatusBadRequest, fmt.Sprintf("Validation failed: %v", err), err)
		utils.RespondWithError(w, appErr)
		return
	}

	user, err := middlewares.GetUserFromContext(r)
	if err != nil {
		appErr := types.NewAppError(http.StatusForbidden, "Forbidden: Unable to retrieve user", err)
		utils.RespondWithError(w, appErr)
		return
	}

	if user.Role != "faculty" && user.Role != "admin" {
		appErr := types.NewAppError(http.StatusUnauthorized, "Unauthorized: Insufficient permissions", nil)
		utils.RespondWithError(w, appErr)
		return
	}

	module := models.Module{
		Title:    createModuleReq.Title,
		CourseID: createModuleReq.CourseID,
	}

	result := db.DB.Create(&module)
	if result.Error != nil {
		appErr := types.NewAppError(http.StatusInternalServerError, "Failed to create module", err)
		utils.RespondWithError(w, appErr)
		return
	}

	types.RespondWithJSON(w, http.StatusCreated, module)
}

func GetModules(w http.ResponseWriter, r *http.Request) {
	courseID := r.URL.Query().Get("course_id")
	if courseID == "" {
		appErr := types.NewAppError(http.StatusBadRequest, "course_id is required", nil)
		utils.RespondWithError(w, appErr)
		return
	}
	var modules []models.Module
	if err := db.DB.Preload("Videos").Where("course_id = ?", courseID).Find(&modules).Error; err != nil {
		appErr := types.NewAppError(http.StatusInternalServerError, "Failed to get modules", err)
		utils.RespondWithError(w, appErr)
		return
	}
	types.RespondWithJSON(w, http.StatusOK, modules)
}
