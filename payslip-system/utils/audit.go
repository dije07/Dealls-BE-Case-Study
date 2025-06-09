package utils

import (
	"time"

	"github.com/dije07/payslip-system/database"
	"github.com/dije07/payslip-system/models"
	"github.com/google/uuid"
)

func LogAudit(action, entity string, entityID, performedBy uuid.UUID, ip, requestID string) {
	log := models.AuditLog{
		Action:      action,
		Entity:      entity,
		EntityID:    entityID,
		PerformedBy: performedBy,
		IP:          ip,
		RequestID:   requestID,
		CreatedAt:   time.Now(),
	}
	_ = database.DB.Create(&log)
}
