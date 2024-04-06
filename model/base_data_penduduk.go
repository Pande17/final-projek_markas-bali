package model

import (
	"time"
)

type BaseData struct {
	id         string    `json:"_id"`
	name       string    `json:"name"`
	email      string    `json:"email"`
	phone      int       `json:"phone"`
	address    string    `json:"address"`
	city       string    `json:"city"`
	postalcode int       `json:"postal_code"`
	Tanggal    time.Time `json:"tanggal"`
}
