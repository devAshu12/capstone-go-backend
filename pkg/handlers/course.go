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

func CreateCourse(w http.ResponseWriter, r *http.Request) {

	var createCourseReq types.CreateCourseReq
	// Decode the request body into the struct
	err := json.NewDecoder(r.Body).Decode(&createCourseReq)
	if err != nil {
		appErr := types.NewAppError(http.StatusBadRequest, "Invalid request format", err)
		utils.RespondWithError(w, appErr)
		return
	}

	// Validate the input request
	if err := config.ValidateRequest(createCourseReq); err != nil {
		appErr := types.NewAppError(http.StatusBadRequest, fmt.Sprintf("Validation failed: %v", err), err)
		utils.RespondWithError(w, appErr)
		return
	}

	// Retrieve the user from the context (ensuring authorization and role checking)
	user, err := middlewares.GetUserFromContext(r)
	if err != nil {
		appErr := types.NewAppError(http.StatusForbidden, "Forbidden: Unable to retrieve user", err)
		utils.RespondWithError(w, appErr)
		return
	}

	// Ensure the user has the role to create a course (e.g., is a faculty or admin)
	if user.Role != "faculty" && user.Role != "admin" {
		appErr := types.NewAppError(http.StatusUnauthorized, "Unauthorized: Insufficient permissions", nil)
		utils.RespondWithError(w, appErr)
		return
	}

	course := models.Course{
		Title:     createCourseReq.Title,
		Price:     createCourseReq.Price,
		UserRefer: user.ID,
	}

	result := db.DB.Create(&course)
	if result.Error != nil {
		appErr := types.NewAppError(http.StatusInternalServerError, "Failed to create course", err)
		utils.RespondWithError(w, appErr)
		return
	}

	types.RespondWithJSON(w, http.StatusCreated, course)
}

func GetCourses(w http.ResponseWriter, r *http.Request) {
	user, err := middlewares.GetUserFromContext(r)
	if err != nil {
		fmt.Print(err)
	}
	fmt.Print(user)
}
