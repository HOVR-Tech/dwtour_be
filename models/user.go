package models

type User struct {
	ID          int           `json:"id"`
	Name        string        `json:"name" gorm:"type: varchar(255)"`
	Gender      string        `json:"gender" gorm:"type: varchar(255)"`
	Email       string        `json:"email" gorm:"type: varchar(255)"`
	Password    string        `json:"password"`
	Number      string        `json:"number" gorm:"type: varchar(255)"`
	Address     string        `json:"address" gorm:"type: varchar(255)"`
	Transaction []Transaction `json:"transaction"`
	Role        string        `json:"role" gorm:"type: varchar(255)"`
}
