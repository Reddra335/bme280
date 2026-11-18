package dto

import (
	"time"
)

type DataRespone struct {
	DeviceId    string    `json:"device_id"`
	TekananKpa  float64   `json:"tekanan_kpa"`
	WaktuDibaca time.Time `json:"waktu_dibaca"`
}

// id BIGINT
// device_id VARCHAR(50)
// tekanan_kpa DOUBLE PRECISION
// waktu_dibaca TIMESTAMPTZ
