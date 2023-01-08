package models

import "time"

type Transaction struct {
	ID        int       `json:"id"`
	Qty       int       `json:"qty" gorm:"type: int"`
	Total     int       `json:"total" gorm:"type: int"`
	Status    string    `json:"status" gorm:"type: varchar(255)"`
	Image     string    `json:"attachment" gorm:"type: varchar(255)"`
	TripID    int       `json:"trip_id" `
	Trip      Trip      `json:"trip"`
	UserID    int       `json:"user_id" `
	User      User      `json:"user"`
	CreatedAt time.Time `json:"-"`
	UpdatedAt time.Time `json:"-"`
}
