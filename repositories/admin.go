package repositories

import (
	"dumbmerch/models"

	"gorm.io/gorm"
)

type AdminRepository interface {
	AddAdmin(admin models.Admin) (models.Admin, error)
}

func RepositoryAdmin(db *gorm.DB) *repository {
	return &repository{db}
}

func (r *repository) AddAdmin(admin models.Admin) (models.Admin, error) {
	err := r.db.Create(&admin).Error

	return admin, err
}
