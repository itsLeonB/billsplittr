package util

import "gorm.io/gorm"

func DefaultOrder() func(*gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Order("created_at DESC")
	}
}
