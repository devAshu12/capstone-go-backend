package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"github/devAshu12/learning_platform_GO_backend/internal/auth"
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
		http.Error(w, "Invalid request format", http.StatusBadRequest)
		return
	}

	// check if email already exist
	var existingUser models.User
	if err := db.DB.Where("email = ?", registerReq.Email).First(&existingUser).Error; err == nil {
		http.Error(w, "Email already in use", http.StatusConflict)
		return
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	// encrypt password
	hashedPassword, err := auth.HashPassword(registerReq.Password)
	if err != nil {
		http.Error(w, "Password encryption failed", http.StatusInternalServerError)
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
		http.Error(w, "Invalid role", http.StatusBadRequest)
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
		http.Error(w, "Failed to create user", http.StatusInternalServerError)
		return
	}

	// create JWT
	access_token, refresh_token, err := auth.GenerateToken(user.ID)
	if err != nil {
		http.Error(w, "Failed to create JWT", http.StatusInternalServerError)
		return
	}

	// set cookies
	http.SetCookie(w, &http.Cookie{
		Name:     "access_token",
		Value:    access_token,
		Expires:  time.Now().Add(24 * time.Hour),
		HttpOnly: true,
		Secure:   true,
	})

	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    refresh_token,
		Expires:  time.Now().Add(14 * 24 * time.Hour),
		HttpOnly: true,
		Secure:   true,
	})

	// send response
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "User registered successfully",
	})
}
