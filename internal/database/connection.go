package database

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func NewConnection(dsn string) (*gorm.DB, error) {
	return gorm.Open(mysql.Open(dsn), &gorm.Config{PrepareStmt: true})
}
