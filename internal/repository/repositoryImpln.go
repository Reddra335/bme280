package repository

import (
	"context"
	"database/sql"
	"menghitung-tekanan/internal/domain"
)

type TekananImpln struct {
	db *sql.DB
}

func NewTekananImpln(db *sql.DB) Tekanan {
	return &TekananImpln{
		db: db,
	}
}

func (t *TekananImpln) Save(ctx context.Context, data domain.Data) error {
	query := `
		INSERT INTO sensor.data_sensor (device_id, temperatur, kelembapan, tekanan_udara, data_dibuat)
		VALUES ($1, $2, $3, $4, $5);
	`

	_, err := t.db.ExecContext(
		ctx,
		query,
		data.DeviceId,
		data.Temperatur,
		data.Kelembapan,
		data.TekananUdara,
		data.DataDibuat,
	)
	return err
}

func (t *TekananImpln) GetAll(ctx context.Context) ([]domain.Data, error) {
	query := `
		SELECT id, device_id, temperatur, kelembapan, tekanan_udara, data_dibuat
		FROM sensor.data_sensor
		ORDER BY data_dibuat DESC
		LIMIT 200
	`

	rows, err := t.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []domain.Data

	for rows.Next() {
		var data domain.Data
		err := rows.Scan(
			&data.ID,
			&data.DeviceId,
			&data.Temperatur,
			&data.Kelembapan,
			&data.TekananUdara,
			&data.DataDibuat,
		)
		if err != nil {
			return nil, err
		}

		result = append(result, data)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return result, nil
}

func (t *TekananImpln) GetLatest(ctx context.Context) (domain.Data, error) {
	query := `
		SELECT id, device_id, temperatur, kelembapan, tekanan_udara, data_dibuat
		FROM sensor.data_sensor
		ORDER BY data_dibuat DESC
		LIMIT 1
	`

	var data domain.Data
	err := t.db.QueryRowContext(ctx, query).Scan(
		&data.ID,
		&data.DeviceId,
		&data.Temperatur,
		&data.Kelembapan,
		&data.TekananUdara,
		&data.DataDibuat,
	)
	if err != nil {
		return domain.Data{}, err
	}

	return data, nil
}
