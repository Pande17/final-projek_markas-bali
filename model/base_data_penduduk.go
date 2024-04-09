package model

import (
	"time"
)

type BaseData struct {
	ID         string    `json:"_id"`
	Name       string    `json:"name"`
	Email      string    `json:"email"`
	Phone      int       `json:"phone"`
	Addres     string    `json:"address"`
	City       string    `json:"city"`
	Postalcode int       `json:"postal_code"`
	Tanggal    time.Time `json:"tanggal"`
}
