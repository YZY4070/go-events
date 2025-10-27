package models

type User struct {
	ID    int64  `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

type Event struct {
	ID          int64  `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
}

type Booking struct {
	ID      int64 `json:"id"`
	EventID int64 `json:"event_id"`
	UserID  int64 `json:"user_id"`
}
