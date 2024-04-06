package model

import (
	"time"
)

type BaseData struct {
	id         string
	name       string
	email      string
	phone      int
	address    string
	city       string
	postalcode int
	Tanggal    time.Time
}
