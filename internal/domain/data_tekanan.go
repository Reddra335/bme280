package domain

import "time"

type Data struct {
	ID           int       `json:"id"`
	DeviceId     string    `json:"device_id"`
	Temperatur   float64   `json:"temperatur"`
	Kelembapan   float64   `json:"kelembapan"`
	TekananUdara float64   `json:"tekananUdara"`
	DataDibuat   time.Time `json:"dataDibuat"`
}
