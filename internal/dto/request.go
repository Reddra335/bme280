package dto

type DataReq struct {
	DeviceId   string  `json:"device_id"`
	TekananKpa float64 `json:"tekanan_kpa"`
}

// id BIGINT
// device_id VARCHAR(50)
// tekanan_kpa DOUBLE PRECISION
// waktu_dibaca TIMESTAMPTZ
