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

func (r *MaterialRepository) Save(m *models.Material) (*models.Material, error) {
	return m, r.db.Save(m).Error
}

func (r *MaterialRepository) FindAll() ([]*models.Material, error) {
	var materials []*models.Material
	return materials, r.db.Find(&materials).Error
}

func (r *MaterialRepository) FindById(id int) (*models.Material, error) {
	var m models.Material
	if err := r.db.First(&m, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &m, nil
}

func (r *MaterialRepository) Delete(m *models.Material) error {
	return r.db.Delete(m).Error
}

func (r *MaterialRepository) FindScarce(threshold float64) ([]*models.Material, error) {
	var materials []*models.Material
	err := r.db.Where("quantity <= ?", threshold).Find(&materials).Error
	return materials, err
}
