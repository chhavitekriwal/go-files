
package models

import (
    "gorm.io/gorm"
    "time"
)

type FileTransaction struct {
    gorm.Model
    UserID   uint		`gorm:"column:user_id"`
    Filename string		`gorm:"column:filename"`
    Type     string 	`gorm:"column:transaction_type"`
    Timestamp time.Time	`gorm:"column:timestamp"`
}