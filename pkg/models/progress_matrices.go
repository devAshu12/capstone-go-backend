package models

import "gorm.io/gorm"

// UserCourseProgress tracks overall course progress
type UserCourseProgress struct {
	gorm.Model
	UserID          string  `gorm:"type:uuid;index" json:"user_id"`
	User            User    `gorm:"foreignKey:UserID" json:"user"`
	CourseID        string  `gorm:"type:uuid;index" json:"course_id"`
	Course          Course  `gorm:"foreignKey:CourseID" json:"course"`
	OverallProgress float64 `gorm:"type:float;not null;default:0" json:"overall_progress"` // 0-100
	IsCompleted     bool    `gorm:"type:bool;default:false" json:"is_completed"`
}

// ModuleProgress tracks progress for individual modules
type ModuleProgress struct {
	gorm.Model
	UserID      string  `gorm:"type:uuid;index" json:"user_id"`
	User        User    `gorm:"foreignKey:UserID" json:"user"`
	ModuleID    string  `gorm:"type:uuid;index" json:"module_id"`
	Module      Module  `gorm:"foreignKey:ModuleID" json:"module"`
	CourseID    string  `gorm:"type:uuid;index" json:"course_id"`
	Course      Course  `gorm:"foreignKey:CourseID" json:"course"`
	Progress    float64 `gorm:"type:float;not null;default:0" json:"progress"` // 0-100
	IsCompleted bool    `gorm:"type:bool;default:false" json:"is_completed"`
}

// QuizScore tracks quiz attempts and scores
type QuizScore struct {
	gorm.Model
	UserID   string  `gorm:"type:uuid;index" json:"user_id"`
	User     User    `gorm:"foreignKey:UserID" json:"user"`
	QuizID   string  `gorm:"type:uuid;index" json:"quiz_id"`
	Quiz     Quiz    `gorm:"foreignKey:QuizID" json:"quiz"`
	Score    float64 `gorm:"type:float;not null" json:"score"` // 0-100
	Attempts int     `gorm:"type:int;not null;default:1" json:"attempts"`
	Passed   bool    `gorm:"type:bool;default:false" json:"passed"`
}

// AssignmentSubmission tracks assignment submissions and grades
type AssignmentSubmission struct {
	gorm.Model
	UserID        string     `gorm:"type:uuid;index" json:"user_id"`
	User          User       `gorm:"foreignKey:UserID" json:"user"`
	AssignmentID  string     `gorm:"type:uuid;index" json:"assignment_id"`
	Assignment    Assignment `gorm:"foreignKey:AssignmentID" json:"assignment"`
	SubmissionURL string     `gorm:"type:text" json:"submission_url"`
	Grade         float64    `gorm:"type:float" json:"grade"`        // 0-100
	Status        string     `gorm:"type:varchar(20)" json:"status"` // pending, graded, resubmit
	Feedback      string     `gorm:"type:text" json:"feedback"`
}
