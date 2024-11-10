package types

type VideoProgress struct {
	UserID     int     `json:"user_id"`
	VideoID    int     `json:"video_id"`
	Progress   float64 `json:"progress"`
	TimeSpent  float32 `json:"time_spent"`
	Completion bool    `json:"completion"`
}

type UserRegisterReq struct {
	FirstName  string `json:"first_name" binding:"required"`
	SecondName string `json:"second_name" binding:"required"`
	Email      string `json:"email" binding:"required,email"`
	Password   string `json:"password" binding:"required,min=8"`
	Role       string `json:"role" binding:"required"` // E.g., "student", "faculty", etc.
}
