package types

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
