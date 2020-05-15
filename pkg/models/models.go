package models

import (
	"errors"
	"time"
)

var ErrNoRecord = errors.New("models: no matching record found")

type Category struct {
	ID      int
	Name    string
	Created time.Time
}

type Product struct {
	ID          int
	CatID       int
	Name        string
	Price       float64
	Stock       int
	ImgURL      string
	Description string
	Created     time.Time
}
