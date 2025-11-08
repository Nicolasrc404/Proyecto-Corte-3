package repository

import (
	"backend-avanzada/models"

	"gorm.io/gorm"
)

type MissionRepository struct {
	db *gorm.DB
}

func NewMissionRepository(db *gorm.DB) *MissionRepository {
	return &MissionRepository{db: db}
}

func (r *MissionRepository) FindAll() ([]*models.Mission, error) {
	var missions []*models.Mission
	err := r.db.Preload("Alchemist.Person").Find(&missions).Error
	return missions, err
}

func (r *MissionRepository) FindById(id int) (*models.Mission, error) {
	var m models.Mission
	err := r.db.Preload("Alchemist.Person").First(&m, id).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return &m, err
}

func (r *MissionRepository) Save(m *models.Mission) (*models.Mission, error) {
	if err := r.db.Save(m).Error; err != nil {
		return nil, err
	}
	return m, nil
}

func (r *MissionRepository) Delete(m *models.Mission) error {
	return r.db.Delete(m).Error
}
