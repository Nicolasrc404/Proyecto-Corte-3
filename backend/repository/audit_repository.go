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
	var audits []*models.Audit
	return audits, r.db.Find(&audits).Error
}
