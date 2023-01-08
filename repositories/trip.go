package repositories

import (
	"dumbmerch/models"

	"gorm.io/gorm"
)

type TripRepository interface {
	FindTrips() ([]models.Trip, error)
	GetTrip(ID int) (models.Trip, error)
	AddTrip(trip models.Trip) (models.Trip, error)
	EditTrip(trip models.Trip, ID int) (models.Trip, error)
	DeleteTrip(trip models.Trip, ID int) (models.Trip, error)
}

func RepositoryTrip(db *gorm.DB) *repository {
	return &repository{db}
}

func (r *repository) FindTrips() ([]models.Trip, error) {
	var trips []models.Trip
	err := r.db.Preload("Country").Find(&trips).Error

	return trips, err
}

func (r *repository) GetTrip(ID int) (models.Trip, error) {
	var trip models.Trip
	err := r.db.Preload("Country").First(&trip, ID).Error

	return trip, err
}

func (r *repository) AddTrip(trip models.Trip) (models.Trip, error) {
	err := r.db.Preload("Country").Create(&trip).Error

	return trip, err
} 

func (r *repository) EditTrip(trip models.Trip, ID int) (models.Trip, error) {
	err := r.db.Model(&trip).Updates(trip).Error
	
	// err := r.db.Raw("UPDATE trips SET title=?, country_id=?, accomodation=?, transportation=?,  eat=?, day=?, night=?, date_trip=?, price=?, quota=?, description=?, image=? WHERE id=?", trip.Title, trip.CountryID, trip.Accomodation, trip.Transportation, trip.Eat,  trip.Day, trip.Night, trip.Date_Trip, trip.Price, trip.Quota, trip.Description, trip.Image, trip.ID).Scan(&trip).Error

	return trip, err
}

func (r *repository) DeleteTrip(trip models.Trip, ID int) (models.Trip, error) {
	err := r.db.Preload("Country").Delete(&trip).Error

	return trip, err
}