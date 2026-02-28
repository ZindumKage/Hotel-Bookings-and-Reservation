package models

import "gorm.io/gorm"

type AuditLog struct {
	gorm.Model
	UserID     uint   `gorm:"index"`
	Action     string `gorm:"index"`
	Entity     string `gorm:"index"`
	RiskLevel  string `gorm:"index"`
	Suspicious bool   `gorm:"index"`
}
