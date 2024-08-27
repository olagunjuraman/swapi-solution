package models

import "time"

type Model struct {
	ID        uint       `gorm:"primary_key" json:"id,omitempty"`
	CreatedAt time.Time  `json:"createdAt"`
	UpdatedAt time.Time  `json:"-" `
	DeletedAt *time.Time `sql:"index" json:"deletedAt,omitempty"`
}

type Comments struct {
	Model
	Comment   string `gorm:"type:text;" json:"comment,omitempty"`
	FilmID    int    `gorm:"type:int;"  json:"film_id,omitempty"`
	Ipaddress string `gorm:"type:varchar(100);"  json:"ipaddress"`
}
