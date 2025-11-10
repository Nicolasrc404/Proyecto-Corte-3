package repository

import (
	"backend-avanzada/models"
	"time"

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

func (r *MissionRepository) FindOpenBefore(threshold time.Time) ([]*models.Mission, error) {
	var ms []*models.Mission
	err := r.db.Where("status <> ? AND updated_at < ?", "completada", threshold).Find(&ms).Error
	return ms, err
}
