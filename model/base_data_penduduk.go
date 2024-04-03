package model

import (
	"time"
)

type BaseData struct {
	ID 		string 
	Name 	string
	Age 	int
	Phone 	int
	Email 	string	
	Tanggal	time.Time
}