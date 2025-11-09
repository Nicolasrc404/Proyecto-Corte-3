package repository

import (
	"backend-avanzada/models"

	"gorm.io/gorm"
)

type MissionRepository struct{ db *gorm.DB }

func NewMissionRepository(db *gorm.DB) *MissionRepository { return &MissionRepository{db: db} }

func (r *MissionRepository) Save(m *models.Mission) (*models.Mission, error) {
	return m, r.db.Save(m).Error
}

func (r *MissionRepository) FindAll() ([]*models.Mission, error) {
	var xs []*models.Mission
	return xs, r.db.Find(&xs).Error
}

func (r *MissionRepository) FindById(id int) (*models.Mission, error) {
	var m models.Mission
	if err := r.db.First(&m, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &m, nil
}

func (r *MissionRepository) Delete(m *models.Mission) error {
	return r.db.Delete(m).Error
}
