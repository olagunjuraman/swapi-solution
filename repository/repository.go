package repository

import (
	"github.com/jinzhu/gorm"
)

type Repository struct {
	DB *gorm.DB
}

func New(db *gorm.DB) *Repository {
	return &Repository{db}
}
