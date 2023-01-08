package repositories

import (
	"dumbmerch/models"

	"gorm.io/gorm"
)

type TransactionRepository interface {
	BookTransaction(transaction models.Transaction) (models.Transaction, error)
	AddTransaction(transaction models.Transaction) (models.Transaction, error)
	FindTransactions() ([]models.Transaction, error)
	GetTransaction(ID int) (models.Transaction, error)
	GetTransactionByUserID(ID int) ([]models.Transaction, error)
	GetOneTransaction(ID string) (models.Transaction, error)
	UpdateTransaction(transaction models.Transaction, ID int) (models.Transaction, error)
	Notification(status string, ID int) (models.Transaction, error)

	// ADMIN
	CheckTransaction(transaction models.Transaction, ID int) (models.Transaction, error)
	DeleteTransaction(transaction models.Transaction, ID int) (models.Transaction, error)
}

func RepositoryTransaction(db *gorm.DB) *repository {
	return &repository{db}
}

func (r *repository) BookTransaction(transaction models.Transaction) (models.Transaction, error) {
	err := r.db.Preload("Trip").Preload("Trip.Country").Preload("User").Create(&transaction).Error

	return transaction, err
}

func (r *repository) AddTransaction(transaction models.Transaction) (models.Transaction, error) {
	err := r.db.Preload("Trip").Preload("Trip.Country").Preload("User").Create(&transaction).Error

	return transaction, err
}

func (r *repository) FindTransactions() ([]models.Transaction, error) {
	var transactions []models.Transaction
	err := r.db.Preload("Trip").Preload("Trip.Country").Preload("User").Find(&transactions).Error

	return transactions, err
}

func (r *repository) GetTransaction(ID int) (models.Transaction, error) {
	var transaction models.Transaction
	err := r.db.Preload("Trip").Preload("Trip.Country").Preload("User").First(&transaction, ID).Error

	return transaction, err
}

func (r *repository) GetTransactionByUserID(ID int) ([]models.Transaction, error) {
	var transaction []models.Transaction
	err := r.db.Preload("Trip").Preload("Trip.Country").Preload("User").Where("user_id =?", ID).Find(&transaction).Error

	return transaction, err
}

func (r *repository) GetOneTransaction(ID string) (models.Transaction, error) {
	var transactions models.Transaction
	err := r.db.Preload("Trip").Preload("Trip.Country").Preload("User").First(&transactions, ID).Error

	return transactions, err
}

func (r *repository) Notification(status string, ID int) (models.Transaction, error) {
	var transaction models.Transaction
	r.db.Debug().Preload("Trip").First(&transaction, ID)

	if status != transaction.Status && status == "success" {
		var trip models.Trip
		r.db.First(&trip, transaction.Trip.ID)
		trip.Quota = trip.Quota - transaction.Qty
		r.db.Model(&trip).Updates(trip)
	}

	transaction.Status = status

	err := r.db.Preload("Trip").Preload("Trip.Country").Preload("User").Model(&transaction).Updates(transaction).Error

	return transaction, err
}

func (r *repository) UpdateTransaction(transaction models.Transaction, ID int) (models.Transaction, error) {
	err := r.db.Preload("Trip").Preload("Trip.Country").Preload("User").Model(&transaction).Updates(transaction).Error

	return transaction, err
}

// ADMIN
func (r *repository) CheckTransaction(transaction models.Transaction, ID int) (models.Transaction, error) {
	err := r.db.Preload("Trip").Preload("Trip.Country").Preload("User").Model(&transaction).Updates(transaction).Error

	return transaction, err
}

func (r *repository) DeleteTransaction(transaction models.Transaction, ID int) (models.Transaction, error) {
	err := r.db.Delete(&transaction, ID).Error

	return transaction, err
}
