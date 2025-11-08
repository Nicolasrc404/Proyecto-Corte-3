package repository

import (
	"backend-avanzada/models"

	"gorm.io/gorm"
)

type MaterialRepository struct {
	db *gorm.DB
}

func NewMaterialRepository(db *gorm.DB) *MaterialRepository {
	return &MaterialRepository{db: db}
}

func (r *MaterialRepository) FindAll() ([]*models.Material, error) {
	var materials []*models.Material
	err := r.db.Find(&materials).Error
	return materials, err
}

func (r *MaterialRepository) FindById(id int) (*models.Material, error) {
	var m models.Material
	err := r.db.First(&m, id).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return &m, err
}

func (r *MaterialRepository) Save(m *models.Material) (*models.Material, error) {
	if err := r.db.Save(m).Error; err != nil {
		return nil, err
	}
	return m, nil
}

func (r *MaterialRepository) Delete(m *models.Material) error {
	return r.db.Delete(m).Error
}
