package handlers

import (
	"github/devAshu12/learning_platform_GO_backend/internal/utils"
	"github/devAshu12/learning_platform_GO_backend/pkg/db"
	"github/devAshu12/learning_platform_GO_backend/pkg/models"
	"github/devAshu12/learning_platform_GO_backend/pkg/types"
	"net/http"
)

// Add these helper functions at the top
func respondWithError(w http.ResponseWriter, status int, message string, err error) {
	appErr := types.NewAppError(status, message, err)
	utils.RespondWithError(w, appErr)
}

func getQueryParams(r *http.Request, params ...string) map[string]string {
	result := make(map[string]string)
	for _, param := range params {
		result[param] = r.URL.Query().Get(param)
	}
	return result
}

// see student details with count of students enrolled in the course
func GetCourseEnrollmentCount(w http.ResponseWriter, r *http.Request) {
	params := getQueryParams(r, "course_id")

	var students []models.User
	if err := db.DB.Model(&models.Course{}).Where("id = ?", params["course_id"]).Association("Students").Find(&students); err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to get students enrolled in the course", err)
		return
	}

	types.RespondWithJSON(w, http.StatusOK, students)
}

// see the progress of the student in the course
func GetProgress(w http.ResponseWriter, r *http.Request) {
	params := getQueryParams(r, "course_id", "student_id", "quiz_id", "type")

	switch params["type"] {
	case "course":
		var progress models.UserCourseProgress
		if err := db.DB.First(&progress, "course_id = ? AND user_id = ?", params["course_id"], params["student_id"]).Error; err != nil {
			respondWithError(w, http.StatusNotFound, "Course progress not found", err)
			return
		}
		types.RespondWithJSON(w, http.StatusOK, progress)

	case "quiz":
		var progress models.QuizScore
		if err := db.DB.First(&progress, "quiz_id = ? AND user_id = ?", params["quiz_id"], params["student_id"]).Error; err != nil {
			respondWithError(w, http.StatusNotFound, "Quiz progress not found", err)
			return
		}
		types.RespondWithJSON(w, http.StatusOK, progress)
	}
}

// get course completion percentage
func GetCourseCompletionPercentage(w http.ResponseWriter, r *http.Request) {
	params := getQueryParams(r, "course_id")

	var stats struct {
		TotalEnrolled  int     `json:"total_enrolled"`
		TotalCompleted int     `json:"total_completed"`
		CompletionRate float64 `json:"completion_rate"`
	}

	// Get enrolled students and completion count in a single query
	err := db.DB.Raw(`
		SELECT 
			COUNT(DISTINCT u.id) as total_enrolled,
			COUNT(DISTINCT CASE WHEN ucp.completion_percentage = 100 THEN u.id END) as total_completed
		FROM users u
		JOIN user_course_progress ucp ON u.id = ucp.user_id
		WHERE ucp.course_id = ?
	`, params["course_id"]).Scan(&stats).Error

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to fetch course statistics", err)
		return
	}

	if stats.TotalEnrolled > 0 {
		stats.CompletionRate = float64(stats.TotalCompleted) / float64(stats.TotalEnrolled) * 100
	}

	types.RespondWithJSON(w, http.StatusOK, stats)
}

// get average score of the quiz
func GetAverageQuizScore(w http.ResponseWriter, r *http.Request) {
	quizID := r.URL.Query().Get("quiz_id")

	// Get all quiz scores for this quiz
	var quizScores []models.QuizScore
	err := db.DB.Where("quiz_id = ?", quizID).Find(&quizScores).Error
	if err != nil {
		appErr := types.NewAppError(http.StatusInternalServerError, "Failed to fetch quiz scores", err)
		utils.RespondWithError(w, appErr)
		return
	}

	// Check if any students attempted the quiz
	if len(quizScores) == 0 {
		types.RespondWithJSON(w, http.StatusOK, map[string]interface{}{
			"average_score": 0,
			"attempts":      0,
		})
		return
	}

	// Calculate total score
	var totalScore float64
	for _, score := range quizScores {
		totalScore += score.Score
	}

	// Calculate average
	averageScore := totalScore / float64(len(quizScores))

	types.RespondWithJSON(w, http.StatusOK, map[string]interface{}{
		"average_score": averageScore,
		"attempts":      len(quizScores),
	})
}

// Get assignment submission statistics for a course
func GetAssignmentStats(w http.ResponseWriter, r *http.Request) {
	params := getQueryParams(r, "course_id", "assignment_id")

	var submissions []models.AssignmentSubmission
	err := db.DB.Where("assignment_id = ?", params["assignment_id"]).Find(&submissions).Error
	if err != nil {
		appErr := types.NewAppError(http.StatusInternalServerError, "Failed to fetch submissions", err)
		utils.RespondWithError(w, appErr)
		return
	}

	// Calculate statistics
	var totalGrade float64
	statusCounts := map[string]int{
		"pending":  0,
		"graded":   0,
		"resubmit": 0,
	}

	for _, submission := range submissions {
		totalGrade += submission.Grade
		statusCounts[submission.Status]++
	}

	avgGrade := 0.0
	if len(submissions) > 0 {
		avgGrade = totalGrade / float64(len(submissions))
	}

	types.RespondWithJSON(w, http.StatusOK, map[string]interface{}{
		"total_submissions": len(submissions),
		"average_grade":     avgGrade,
		"status_breakdown":  statusCounts,
	})
}

// Get module completion statistics
func GetModuleStats(w http.ResponseWriter, r *http.Request) {
	params := getQueryParams(r, "module_id")

	var moduleProgress []models.ModuleProgress
	err := db.DB.Where("module_id = ?", params["module_id"]).Find(&moduleProgress).Error
	if err != nil {
		appErr := types.NewAppError(http.StatusInternalServerError, "Failed to fetch module progress", err)
		utils.RespondWithError(w, appErr)
		return
	}

	var totalProgress float64
	completedCount := 0
	for _, progress := range moduleProgress {
		totalProgress += progress.Progress
		if progress.IsCompleted {
			completedCount++
		}
	}

	avgProgress := 0.0
	if len(moduleProgress) > 0 {
		avgProgress = totalProgress / float64(len(moduleProgress))
	}

	types.RespondWithJSON(w, http.StatusOK, map[string]interface{}{
		"total_students":   len(moduleProgress),
		"completed_count":  completedCount,
		"average_progress": avgProgress,
		"completion_rate":  float64(completedCount) / float64(len(moduleProgress)) * 100,
	})
}

// Get struggling students (low progress or failed attempts)
func GetStrugglingStudents(w http.ResponseWriter, r *http.Request) {
	params := getQueryParams(r, "course_id")
	progressThreshold := 30.0 // Consider students below 30% progress as struggling

	var strugglingStudents []struct {
		models.User
		Progress      float64
		FailedQuizzes int
	}

	// Find students with low progress or multiple failed quizzes
	err := db.DB.Table("users").
		Select("users.*, user_course_progress.overall_progress as progress, "+
			"COUNT(CASE WHEN quiz_scores.passed = false THEN 1 END) as failed_quizzes").
		Joins("JOIN user_course_progress ON users.id = user_course_progress.user_id").
		Joins("LEFT JOIN quiz_scores ON users.id = quiz_scores.user_id").
		Where("user_course_progress.course_id = ? AND "+
			"(user_course_progress.overall_progress < ? OR quiz_scores.passed = false)",
			params["course_id"], progressThreshold).
		Group("users.id, user_course_progress.overall_progress").
		Find(&strugglingStudents).Error

	if err != nil {
		appErr := types.NewAppError(http.StatusInternalServerError, "Failed to fetch struggling students", err)
		utils.RespondWithError(w, appErr)
		return
	}

	types.RespondWithJSON(w, http.StatusOK, strugglingStudents)
}

//see the progress of the student in the assignment

//see the progress of the student in the video

//see the progress of the student in the module

//see the progress of the student in the course
