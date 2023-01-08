package models

type Trip struct {
	ID        int               `json:"id" gorm:"primary_key:auto_increment"`
	Title     string            `json:"title" gorm:"type: varchar(255)"`
	CountryID int               `json:"country_id"`
	Country   CountriesResponse `json:"country"`
	// UserId         int               `json:"user_id"`
	// User           User              `json:"user"`
	Accomodation   string `json:"accomodation" gorm:"type: varchar(255)"`
	Transportation string `json:"transportation" gorm:"type: varchar(255)"`
	Eat            string `json:"eat" gorm:"type: varchar(255)"`
	Duration       string `json:"duration" gorm:"type: varchar(255)"`
	Date_Trip      string `json:"date_trip" gorm:"type: varchar(255)"`
	Price          int    `json:"price" gorm:"type: int"`
	Quota          int    `json:"quota" gorm:"type: int"`
	Description    string `json:"description" gorm:"type: varchar(255)"`
	Image          string `json:"image" gorm:"type: varchar(255)"`
}
