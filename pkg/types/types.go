package types

type VideoProgress struct {
	UserID     int     `json:"user_id"`
	VideoID    int     `json:"video_id"`
	Progress   float64 `json:"progress"`
	TimeSpent  float32 `json:"time_spent"`
	Completion bool    `json:"completion"`
}
