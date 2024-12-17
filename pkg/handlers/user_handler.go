package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"github/devAshu12/learning_platform_GO_backend/internal/auth"
	"github/devAshu12/learning_platform_GO_backend/internal/utils"
	"github/devAshu12/learning_platform_GO_backend/pkg/config"
	"github/devAshu12/learning_platform_GO_backend/pkg/db"
	"github/devAshu12/learning_platform_GO_backend/pkg/models"
	"github/devAshu12/learning_platform_GO_backend/pkg/types"
	"net/http"
	"sync"
	"time"

	"gorm.io/gorm"
)

var (
	queue      []types.VideoProgress
	queueMutex sync.Mutex
	Quit       chan struct{}
)

func batchInsert(progressBatch []types.VideoProgress) {
	// query := `INSERT INTO user_progress (user_id, video_id, progress, time_spent, completion, updated_at)
	//           VALUES ($1, $2, $3, $4, $5, CURRENT_TIMESTAMP)
	//           ON CONFLICT (user_id, video_id)
	//           DO UPDATE SET progress = $3, time_spent = $4, completion = $5, updated_at = CURRENT_TIMESTAMP`
	for _, progress := range progressBatch {
		fmt.Println("Processing and saving to DB:", progress.VideoID)
	}
}

func ProcessQueue(interval time.Duration, batchSize int) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {

		case <-ticker.C:
			dispatchBatch(batchSize)

		case <-Quit: // Stop processing when quit signal is received
			fmt.Println("Shutting down queue processor...")
			dispatchBatch(0) // Flush remaining items
			return
		}
	}
}

func dispatchBatch(batchSize int) {
	queueMutex.Lock()
	defer queueMutex.Unlock()

	if len(queue) >= batchSize || (batchSize == 0 && len(queue) > 0) {
		progressBatch := make([]types.VideoProgress, len(queue))
		copy(progressBatch, queue)
		queue = []types.VideoProgress{}
		go batchInsert(progressBatch)
	}
}

func UpdateProgress(w http.ResponseWriter, r *http.Request) {
	var progress types.VideoProgress
	err := json.NewDecoder(r.Body).Decode(&progress)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	queueMutex.Lock()
	queue = append(queue, progress)
	fmt.Println("Added to queue:", progress.VideoID)
	queueMutex.Unlock()

	w.WriteHeader(http.StatusAccepted)
}

func Register(w http.ResponseWriter, r *http.Request) {
	var registerReq types.UserRegisterReq
	err := json.NewDecoder(r.Body).Decode(&registerReq)
	if err != nil {
		appErr := types.NewAppError(http.StatusBadRequest, "Invalid request format", err)
		utils.RespondWithError(w, appErr)
		return
	}

	if err := config.ValidateRequest(registerReq); err != nil {
		appErr := types.NewAppError(http.StatusBadRequest, fmt.Sprintf("Validation failed: %v", err), err)
		utils.RespondWithError(w, appErr)
		return
	}

	// check if email already exist
	var existingUser models.User
	if err := db.DB.Where("email = ?", registerReq.Email).First(&existingUser).Error; err == nil {
		appErr := types.NewAppError(http.StatusConflict, "Email already in use", err)
		utils.RespondWithError(w, appErr)
		return
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		appErr := types.NewAppError(http.StatusInternalServerError, "Database error", err)
		utils.RespondWithError(w, appErr)
		return
	}

	// encrypt password
	hashedPassword, err := auth.HashPassword(registerReq.Password)
	if err != nil {
		appErr := types.NewAppError(http.StatusInternalServerError, "Password encryption failed", err)
		utils.RespondWithError(w, appErr)
		return
	}

	var role models.RoleType
	switch registerReq.Role {
	case "super_admin_dev":
		role = models.SuperAdminDev
	case "super_admin":
		role = models.SuperAdmin
	case "faculty":
		role = models.Faculty
	case "student":
		role = models.Student
	default:
		appErr := types.NewAppError(http.StatusBadRequest, "Invalid role", err)
		utils.RespondWithError(w, appErr)
		return
	}

	// create user
	user := models.User{
		FirstName:  registerReq.FirstName,
		SecondName: registerReq.SecondName,
		Email:      registerReq.Email,
		Password:   hashedPassword,
		Role:       role,
	}

	if err := db.DB.Create(&user).Error; err != nil {
		appErr := types.NewAppError(http.StatusInternalServerError, "Failed to create user", err)
		utils.RespondWithError(w, appErr)
		return
	}

	// create JWT
	access_token, refresh_token, err := auth.GenerateToken(user.ID)
	if err != nil {
		fmt.Println(err)
		appErr := types.NewAppError(http.StatusInternalServerError, "Failed to create JWT", err)
		utils.RespondWithError(w, appErr)
		return
	}

	// set cookies
	http.SetCookie(w, &http.Cookie{
		Name:     "access_token",
		Value:    access_token,
		Expires:  time.Now().Add(24 * time.Hour),
		HttpOnly: true,
		Secure:   true,
		Path:     "/",
		SameSite: http.SameSiteStrictMode,
	})

	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    refresh_token,
		Expires:  time.Now().Add(14 * 24 * time.Hour),
		HttpOnly: true,
		Secure:   true,
		Path:     "/",
		SameSite: http.SameSiteStrictMode,
	})

	// send response
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "User registered successfully",
		"success": true,
	})
}

func Login(w http.ResponseWriter, r *http.Request) {

	var userLoginReq types.UserLoginReq

	err := json.NewDecoder(r.Body).Decode(&userLoginReq)
	if err != nil {
		appErr := types.NewAppError(http.StatusBadRequest, "Invalid request format", err)
		utils.RespondWithError(w, appErr)
		return
	}

	if err := config.ValidateRequest(userLoginReq); err != nil {
		appErr := types.NewAppError(http.StatusBadRequest, fmt.Sprintf("Validation failed: %v", err), err)
		utils.RespondWithError(w, appErr)
		return
	}

	var isUserExist models.User
	if err := db.DB.Where("email = ?", userLoginReq.Email).First(&isUserExist).Error; errors.Is(err, gorm.ErrRecordNotFound) {
		appErr := types.NewAppError(http.StatusNotFound, "Invalid user email", err)
		utils.RespondWithError(w, appErr)
		return
	}

	isMatchErr := auth.CheckPassword(isUserExist.Password, userLoginReq.Password)
	if isMatchErr != nil {
		appError := types.NewAppError(http.StatusUnauthorized, "Invalid password", err)
		utils.RespondWithError(w, appError)
		return
	}

	access_token, refresh_token, err := auth.GenerateToken(isUserExist.ID)
	if err != nil {
		fmt.Println(err)
		appErr := types.NewAppError(http.StatusInternalServerError, "Failed to create JWT", err)
		utils.RespondWithError(w, appErr)
		return
	}

	// set cookies
	http.SetCookie(w, &http.Cookie{
		Name:     "access_token",
		Value:    access_token,
		Expires:  time.Now().Add(24 * time.Hour),
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
	})

	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    refresh_token,
		Expires:  time.Now().Add(14 * 24 * time.Hour),
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
	})

	// send response
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "User logged in successfully",
		"success": true,
	})

}
