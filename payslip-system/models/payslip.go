package models

import (
	"time"

	"github.com/google/uuid"
)

type Payslip struct {
	ID             uuid.UUID `gorm:"type:uuid;primaryKey"`
	UserID         uuid.UUID
	PeriodID       uuid.UUID
	BaseSalary     float64
	AttendanceDays int
	OvertimeHours  int
	OvertimePay    float64
	Reimbursement  float64
	TakeHomePay    float64
	IPAddress      string
	CreatedAt      time.Time
	CreatedBy      uuid.UUID `gorm:"type:uuid"`
	UpdatedBy      uuid.UUID `gorm:"type:uuid"`
}
