package models

type Admin struct {
	ID       int    `json:"id"`
	Name     string `json:"name" gorm:"type: varchar(255)"`
	Email    string `json:"email" gorm:"type: varchar(255)"`
	Password string `json:"password"`
	Role     string `json:"role" gorm:"type: varchar(255)"`
}
