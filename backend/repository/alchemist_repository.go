package repository

import (
	"backend-avanzada/models"

	"gorm.io/gorm"
)

type AlchemistRepository struct {
	db *gorm.DB
}

func NewAlchemistRepository(db *gorm.DB) *AlchemistRepository {
	return &AlchemistRepository{db: db}
}

func (r *AlchemistRepository) FindAll() ([]*models.Alchemist, error) {
	var xs []*models.Alchemist
	err := r.db.Find(&xs).Error
	return xs, err
}

func (r *AlchemistRepository) FindById(id int) (*models.Alchemist, error) {
	var a models.Alchemist
	err := r.db.First(&a, id).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return &a, err
}

func (r *AlchemistRepository) Save(a *models.Alchemist) (*models.Alchemist, error) {
	if err := r.db.Save(a).Error; err != nil {
		return nil, err
	}
	return a, nil
}

func (r *AlchemistRepository) Delete(a *models.Alchemist) error {
	return r.db.Delete(a).Error
}
