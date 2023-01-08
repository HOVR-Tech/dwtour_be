package models

type Country struct {
	ID   int    `json:"id"`
	Name string `json:"name" gorm:"type: varchar(255)"`
}

type CountriesResponse struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

func (CountriesResponse) TableName() string {
	return "countries"
}