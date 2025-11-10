package repository

import (
	"backend-avanzada/models"
	"time"

	"gorm.io/gorm"
)

type TransmutationRepository struct {
	db *gorm.DB
}

func (r *TransmutationRepository) FindPendingBefore(threshold time.Time) ([]*models.Transmutation, error) {
	var ts []*models.Transmutation
	err := r.db.Where("status = ? AND created_at < ?", "en_proceso", threshold).Find(&ts).Error
	return ts, err
}

func (r *TransmutationRepository) FindById(id int) (*models.Transmutation, error) {
	var t models.Transmutation
	err := r.db.First(&t, id).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return &t, err
}

func (r *TransmutationRepository) Delete(t *models.Transmutation) error {
	return r.db.Delete(t).Error
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
