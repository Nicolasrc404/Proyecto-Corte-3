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

func (r *AuditRepository) Save(a *models.Audit) (*models.Audit, error) {
	return a, r.db.Save(a).Error
}

func (r *AuditRepository) FindAll() ([]*models.Audit, error) {
	var audits []*models.Audit
	return audits, r.db.Find(&audits).Error
}

func (r *AuditRepository) FindById(id int) (*models.Audit, error) {
	var a models.Audit
	if err := r.db.First(&a, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &a, nil
}

func (r *AuditRepository) Delete(a *models.Audit) error {
	return r.db.Delete(a).Error
}
