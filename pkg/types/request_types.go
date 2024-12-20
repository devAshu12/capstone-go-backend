package types

import (
	"time"
)

type VideoProgress struct {
	UserID     int     `json:"user_id"`
	VideoID    int     `json:"video_id"`
	Progress   float64 `json:"progress"`
	TimeSpent  float32 `json:"time_spent"`
	Completion bool    `json:"completion"`
}

type UserRegisterReq struct {
	FirstName  string `json:"first_name" validate:"required"`
	SecondName string `json:"second_name" validate:"required"`
	Email      string `json:"email" validate:"required,email"`
	Password   string `json:"password" validate:"required,min=8"`
	Role       string `json:"role" validate:"required"` // E.g., "student", "faculty", etc.
}

type UserLoginReq struct {
	Email    string `json:"email" validate:"required"`
	Password string `json:"password" validate:"required,min=8"`
}

type CreateCourseReq struct {
	Title string  `json:"title" validate:"required"`
	Price float64 `json:"price" validate:"required"`
}

type CreateModuleReq struct {
	Title    string `json:"title" validate:"required"`
	CourseID string `json:"course_id" validate:"required"`
}

type CreateAssignmentReq struct {
	Title       string    `json:"title"`
	Description string    `json:"description"`
	ModuleID    string    `json:"module_id"`
	CourseID    string    `json:"course_id"`
	Deadline    time.Time `json:"deadline"`
}

type UpdateAssignmentReq struct {
	ID          string    `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	ModuleID    string    `json:"module_id"`
	CourseID    string    `json:"course_id"`
	Deadline    time.Time `json:"deadline"`
}

type CreateQuizzReq struct {
	Title       string    `json:"title"`
	Description string    `json:"description"`
	DueDate     time.Time `json:"due_date"`
	CourseID    string    `json:"course_id"`
	ModuleID    *string   `json:"module_id"`
	IsFinal     bool      `json:"is_final"`
}

type UpdateQuizReq struct {
	Title       string    `json:"title"`
	Description string    `json:"description"`
	DueDate     time.Time `json:"due_date"`
	IsFinal     bool      `json:"is_final"`
}

type CreateQuestionReq struct {
	QuizID  string        `json:"quiz_id"`
	Text    string        `json:"text"`
	Points  int           `json:"points"`
	Options []OptionInput `json:"options"`
}

type OptionInput struct {
	Text      string `json:"text"`
	IsCorrect bool   `json:"is_correct"`
}

// type CreateVideoReq struct {
// 	Title    string `form:"title" validate:"required"`
// 	ModuleID string `form:"module_id" validate:"required,uuid"`
// 	File     []byte `form:"file" validate:"required"` // File will be read from multipart form
// }
