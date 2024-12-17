package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"github/devAshu12/learning_platform_GO_backend/internal/middlewares"
	"github/devAshu12/learning_platform_GO_backend/internal/utils"
	"github/devAshu12/learning_platform_GO_backend/pkg/db"
	"github/devAshu12/learning_platform_GO_backend/pkg/models"
	"github/devAshu12/learning_platform_GO_backend/pkg/types"
	"net/http"
	"strings"
)

func CreateQuizz(w http.ResponseWriter, r *http.Request) {
	//get request body
	var createQuizzReq types.CreateQuizzReq
	err := json.NewDecoder(r.Body).Decode(&createQuizzReq)
	if err != nil {
		appErr := types.NewAppError(http.StatusBadRequest, "Invalid request format", err)
		utils.RespondWithError(w, appErr)
		return
	}

	// Retrieve the user from the context (ensuring authorization and role checking)
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

	// check if the course exists
	var course models.Course
	if result := db.DB.First(&course, "id = ?", createQuizzReq.CourseID); result.Error != nil {
		appErr := types.NewAppError(http.StatusNotFound, "Course not found", result.Error)
		utils.RespondWithError(w, appErr)
		return
	}

	quiz := models.Quiz{
		Title:       createQuizzReq.Title,
		Description: createQuizzReq.Description,
		DueDate:     createQuizzReq.DueDate,
		CourseID:    createQuizzReq.CourseID,
		ModuleID:    createQuizzReq.ModuleID,
		IsFinal:     createQuizzReq.IsFinal,
	}

	err = db.DB.Create(&quiz).Error
	if err != nil {
		appErr := types.NewAppError(http.StatusInternalServerError, "Failed to create quiz", err)
		utils.RespondWithError(w, appErr)
		return
	}

	types.RespondWithJSON(w, http.StatusOK, "Quiz created successfully")
}

// add question to quiz
func AddQuestionToQuiz(w http.ResponseWriter, r *http.Request) {
	var questionReq []types.CreateQuestionReq
	err := json.NewDecoder(r.Body).Decode(&questionReq)
	if err != nil {
		appErr := types.NewAppError(http.StatusBadRequest, "Invalid request format", err)
		utils.RespondWithError(w, appErr)
		return
	}

	//check if user is admin and faculty
	user, err := middlewares.GetUserFromContext(r)
	if err != nil {
		appErr := types.NewAppError(http.StatusForbidden, "Forbidden: Unable to retrieve user", err)
		utils.RespondWithError(w, appErr)
		return
	}
	if user.Role != "admin" && user.Role != "faculty" {
		appErr := types.NewAppError(http.StatusUnauthorized, "Unauthorized: Insufficient permissions", nil)
		utils.RespondWithError(w, appErr)
		return
	}

	if len(questionReq) == 0 {
		appErr := types.NewAppError(http.StatusBadRequest, "No questions provided", nil)
		utils.RespondWithError(w, appErr)
		return
	}

	var quiz models.Quiz
	if result := db.DB.First(&quiz, "id = ?", questionReq[0].QuizID); result.Error != nil {
		appErr := types.NewAppError(http.StatusNotFound, "Quiz not found", result.Error)
		utils.RespondWithError(w, appErr)
		return
	}

	var errorArr []string

	for _, qReq := range questionReq {
		// Create question
		question := models.Question{
			QuizID: qReq.QuizID,
			Text:   qReq.Text,
			Points: qReq.Points,
		}

		if err := db.DB.Create(&question).Error; err != nil {
			errorArr = append(errorArr, fmt.Sprintf("Failed to create question: %v", err))
			continue
		}

		// Create options for the question
		for _, optReq := range qReq.Options {
			option := models.Option{
				QuestionID: question.ID,
				Text:       optReq.Text,
				IsCorrect:  optReq.IsCorrect,
			}

			if err := db.DB.Create(&option).Error; err != nil {
				errorArr = append(errorArr, fmt.Sprintf("Failed to create option: %v", err))
			}
		}
	}

	if len(errorArr) > 0 {
		errMsg := strings.Join(errorArr, "; ")
		appErr := types.NewAppError(http.StatusInternalServerError, "Failed to create questions", errors.New(errMsg))
		utils.RespondWithError(w, appErr)
		return
	}

	types.RespondWithJSON(w, http.StatusOK, "Questions and options added successfully")
}

// get all quizzes
func GetAllQuizzes(w http.ResponseWriter, r *http.Request) {
	// Get query parameters
	moduleID := r.URL.Query().Get("module_id")
	courseID := r.URL.Query().Get("course_id")
	userID := r.URL.Query().Get("user_id")

	isFinal := r.URL.Query().Get("is_final")

	// Start building the query prelode questions and options
	query := db.DB.Model(&models.Quiz{}).Preload("Questions").Preload("Questions.Options")

	// Add conditions only if the parameters are provided
	if moduleID != "" {
		query = query.Where("module_id = ?", moduleID)
	}
	if courseID != "" {
		query = query.Where("course_id = ?", courseID)
	}
	if userID != "" {
		query = query.Where("user_id = ?", userID)
	}

	if isFinal == "true" {
		query = query.Where("is_final = ?", true)
	}

	// Execute the query
	var quizzes []models.Quiz
	if err := query.Find(&quizzes).Error; err != nil {
		appErr := types.NewAppError(http.StatusInternalServerError, "Failed to fetch quizzes", err)
		utils.RespondWithError(w, appErr)
		return
	}

	types.RespondWithJSON(w, http.StatusOK, quizzes)
}

// get quiz by id
func GetQuizById(w http.ResponseWriter, r *http.Request) {
	quizID := r.URL.Query().Get("quiz_id")
	var quiz models.Quiz
	err := db.DB.Preload("Questions").Preload("Questions.Options").First(&quiz, "id = ?", quizID).Error
	if err != nil {
		appErr := types.NewAppError(http.StatusNotFound, "Quiz not found", err)
		utils.RespondWithError(w, appErr)
		return
	}
	types.RespondWithJSON(w, http.StatusOK, quiz)
}

// update quiz
func UpdateQuiz(w http.ResponseWriter, r *http.Request) {
	quizID := r.URL.Query().Get("quiz_id")

	var updateQuizReq types.UpdateQuizReq
	err := json.NewDecoder(r.Body).Decode(&updateQuizReq)
	if err != nil {
		appErr := types.NewAppError(http.StatusBadRequest, "Invalid request format", err)
		utils.RespondWithError(w, appErr)
		return
	}

	//check if user is admin and faculty
	user, err := middlewares.GetUserFromContext(r)
	if err != nil {
		appErr := types.NewAppError(http.StatusForbidden, "Forbidden: Unable to retrieve user", err)
		utils.RespondWithError(w, appErr)
		return
	}

	if user.Role != "admin" && user.Role != "faculty" {
		appErr := types.NewAppError(http.StatusUnauthorized, "Unauthorized: Insufficient permissions", nil)
		utils.RespondWithError(w, appErr)
		return
	}

	var quiz models.Quiz
	err = db.DB.First(&quiz, "id = ?", quizID).Error
	if err != nil {
		appErr := types.NewAppError(http.StatusNotFound, "Quiz not found", err)
		utils.RespondWithError(w, appErr)
		return
	}

	quiz.Title = updateQuizReq.Title
	quiz.Description = updateQuizReq.Description
	quiz.DueDate = updateQuizReq.DueDate
	quiz.IsFinal = updateQuizReq.IsFinal

	err = db.DB.Save(&quiz).Error
	if err != nil {
		appErr := types.NewAppError(http.StatusInternalServerError, "Failed to update quiz", err)
		utils.RespondWithError(w, appErr)
		return
	}

	types.RespondWithJSON(w, http.StatusOK, quiz)
}

// delete quiz
func DeleteQuiz(w http.ResponseWriter, r *http.Request) {
	quizID := r.URL.Query().Get("quiz_id")

	//check if user is admin and faculty
	user, err := middlewares.GetUserFromContext(r)
	if err != nil {
		appErr := types.NewAppError(http.StatusForbidden, "Forbidden: Unable to retrieve user", err)
		utils.RespondWithError(w, appErr)
		return
	}

	if user.Role != "admin" && user.Role != "faculty" {
		appErr := types.NewAppError(http.StatusUnauthorized, "Unauthorized: Insufficient permissions", nil)
		utils.RespondWithError(w, appErr)
		return
	}

	err = db.DB.Delete(&models.Quiz{}, "id = ?", quizID).Error
	if err != nil {
		appErr := types.NewAppError(http.StatusInternalServerError, "Failed to delete quiz", err)
		utils.RespondWithError(w, appErr)
		return
	}

	types.RespondWithJSON(w, http.StatusOK, "Quiz deleted successfully")
}
