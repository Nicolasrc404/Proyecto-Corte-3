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
	var alchemists []*models.Alchemist
	err := r.db.Preload("Person").Find(&alchemists).Error
	return alchemists, err
}

func (r *AlchemistRepository) FindById(id int) (*models.Alchemist, error) {
	var a models.Alchemist
	err := r.db.Preload("Person").First(&a, id).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return &a, err
}

func (r *AlchemistRepository) Save(a *models.Alchemist) (*models.Alchemist, error) {
	err := r.db.Save(a).Error
	return a, err
}

func (r *AlchemistRepository) Delete(a *models.Alchemist) error {
	return r.db.Delete(a).Error
}
