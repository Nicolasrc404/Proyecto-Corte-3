package repository

import (
	"backend-avanzada/models"

	"gorm.io/gorm"
)

type TransmutationRepository struct {
	db *gorm.DB
}

func NewTransmutationRepository(db *gorm.DB) *TransmutationRepository {
	return &TransmutationRepository{db: db}
}

func (r *TransmutationRepository) FindAll() ([]*models.Transmutation, error) {
	var ts []*models.Transmutation
	err := r.db.Find(&ts).Error
	return ts, err
}

func (r *TransmutationRepository) Save(t *models.Transmutation) (*models.Transmutation, error) {
	if err := r.db.Save(t).Error; err != nil {
		return nil, err
	}
	return t, nil
}
