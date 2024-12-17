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
	"strings"
)

func CreateAssignment(w http.ResponseWriter, r *http.Request) {
	var createAssignmentReq types.CreateAssignmentReq
	// Decode the request body into the struct
	err := json.NewDecoder(r.Body).Decode(&createAssignmentReq)
	if err != nil {
		appErr := types.NewAppError(http.StatusBadRequest, "Invalid request format", err)
		utils.RespondWithError(w, appErr)
		return
	}

	// Validate the input request
	if err := config.ValidateRequest(createAssignmentReq); err != nil {
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

	assignment := models.Assignment{
		Title:       createAssignmentReq.Title,
		Description: createAssignmentReq.Description,
		ModuleID:    createAssignmentReq.ModuleID,
		CourseID:    createAssignmentReq.CourseID,
		Deadline:    createAssignmentReq.Deadline,
		FacultyID:   user.ID,
	}

	// Verify that the course exists and faculty has access to it
	var course models.Course
	if result := db.DB.First(&course, "id = ?", assignment.CourseID); result.Error != nil {
		appErr := types.NewAppError(http.StatusNotFound, "Course not found", result.Error)
		utils.RespondWithError(w, appErr)
		return
	}

	// Verify that the module belongs to the specified course
	var module models.Module
	if result := db.DB.First(&module, "id = ? AND course_id = ?", assignment.ModuleID, assignment.CourseID); result.Error != nil {
		appErr := types.NewAppError(http.StatusBadRequest, "Module does not belong to the specified course", result.Error)
		utils.RespondWithError(w, appErr)
		return
	}

	result := db.DB.Create(&assignment)
	if result.Error != nil {
		appErr := types.NewAppError(http.StatusInternalServerError, "Failed to create assignment", result.Error)
		utils.RespondWithError(w, appErr)
		return
	}

	types.RespondWithJSON(w, http.StatusCreated, assignment)
}

func GetAssignments(w http.ResponseWriter, r *http.Request) {
	// Parse query parameters
	moduleID := r.URL.Query().Get("module_id")
	courseID := r.URL.Query().Get("course_id")
	facultyID := r.URL.Query().Get("faculty_id")

	var assignments []models.Assignment
	query := db.DB

	// Build conditions slice for all non-empty parameters
	conditions := make([]string, 0)
	args := make([]interface{}, 0)

	if moduleID != "" {
		conditions = append(conditions, "module_id = ?")
		args = append(args, moduleID)
	}
	if courseID != "" {
		conditions = append(conditions, "course_id = ?")
		args = append(args, courseID)
	}
	if facultyID != "" {
		conditions = append(conditions, "faculty_id = ?")
		args = append(args, facultyID)
	}

	// Apply all conditions at once if any exist
	if len(conditions) > 0 {
		query = query.Where(strings.Join(conditions, " AND "), args...)
	}

	// Execute the query
	result := query.Find(&assignments)
	if result.Error != nil {
		appErr := types.NewAppError(http.StatusInternalServerError, "Failed to fetch assignments", result.Error)
		utils.RespondWithError(w, appErr)
		return
	}

	types.RespondWithJSON(w, http.StatusOK, assignments)
}

func GetAssignment(w http.ResponseWriter, r *http.Request) {
	assignmentID := r.URL.Query().Get("assignment_id")

	var assignment models.Assignment
	if result := db.DB.First(&assignment, "id = ?", assignmentID); result.Error != nil {
		appErr := types.NewAppError(http.StatusNotFound, "Assignment not found", result.Error)
		utils.RespondWithError(w, appErr)
		return
	}

	types.RespondWithJSON(w, http.StatusOK, assignment)
}

func UpdateAssignment(w http.ResponseWriter, r *http.Request) {
	var updateAssignmentReq types.UpdateAssignmentReq
	err := json.NewDecoder(r.Body).Decode(&updateAssignmentReq)
	if err != nil {
		appErr := types.NewAppError(http.StatusBadRequest, "Invalid request format", err)
		utils.RespondWithError(w, appErr)
		return
	}

	// check if the assignment exists
	var assignment models.Assignment
	if result := db.DB.First(&assignment, "id = ?", updateAssignmentReq.ID); result.Error != nil {
		appErr := types.NewAppError(http.StatusNotFound, "Assignment not found", result.Error)
		utils.RespondWithError(w, appErr)
		return
	}

	// check if the module exists
	var module models.Module
	if result := db.DB.First(&module, "id = ?", updateAssignmentReq.ModuleID); result.Error != nil {
		appErr := types.NewAppError(http.StatusNotFound, "Module not found", result.Error)
		utils.RespondWithError(w, appErr)
		return
	}

	// check if the course exists
	var course models.Course
	if result := db.DB.First(&course, "id = ?", updateAssignmentReq.CourseID); result.Error != nil {
		appErr := types.NewAppError(http.StatusNotFound, "Course not found", result.Error)
		utils.RespondWithError(w, appErr)
		return
	}

	// update the assignment
	assignment.Title = updateAssignmentReq.Title
	assignment.Description = updateAssignmentReq.Description
	assignment.ModuleID = updateAssignmentReq.ModuleID
	assignment.CourseID = updateAssignmentReq.CourseID
	assignment.Deadline = updateAssignmentReq.Deadline

	db.DB.Save(&assignment)

	types.RespondWithJSON(w, http.StatusOK, assignment)
}

func DeleteAssignment(w http.ResponseWriter, r *http.Request) {
	assignmentID := r.URL.Query().Get("assignment_id")

	var assignment models.Assignment
	if result := db.DB.First(&assignment, "id = ?", assignmentID); result.Error != nil {
		appErr := types.NewAppError(http.StatusNotFound, "Assignment not found", result.Error)
		utils.RespondWithError(w, appErr)
		return
	}

	err := db.DB.Delete(&assignment).Error
	if err != nil {
		appErr := types.NewAppError(http.StatusInternalServerError, "Failed to delete assignment", err)
		utils.RespondWithError(w, appErr)
		return
	}

	types.RespondWithJSON(w, http.StatusOK, "Assignment deleted successfully")
}
