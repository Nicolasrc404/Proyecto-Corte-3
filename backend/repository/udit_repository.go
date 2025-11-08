package repository

import (
	"backend-avanzada/models"

	"gorm.io/gorm"
)

type AuditRepository struct {
	db *gorm.DB
}

func NewAuditRepository(db *gorm.DB) *AuditRepository {
	return &AuditRepository{db: db}
}

func (r *AuditRepository) FindAll() ([]*models.Audit, error) {
	var xs []*models.Audit
	err := r.db.Find(&xs).Error
	return xs, err
}

func (r *AuditRepository) Save(a *models.Audit) (*models.Audit, error) {
	if err := r.db.Save(a).Error; err != nil {
		return nil, err
	}
	return a, nil
}
