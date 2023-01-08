package tripsdto

type TripRequest struct {
	Title          string `json:"title" `
	CountryID      int    `json:"country_id"`
	Accomodation   string `json:"accomodation" `
	Transportation string `json:"transportation" `
	Eat            string `json:"eat" `
	Duration       string `json:"duration" `
	Date_Trip      string `json:"date_trip"`
	Price          int    `json:"price" `
	Quota          int    `json:"quota" `
	Description    string `json:"description"`
	Image          string `json:"image" `
}
