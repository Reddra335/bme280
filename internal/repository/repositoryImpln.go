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
	query := `INSERT INTO tekanan.data (device_id, tekanan_kpa, waktu_dibaca) VALUES ($1,$2,$3);`

	_, err := t.db.ExecContext(ctx, query, data.DeviceId, data.TekananKpa, data.WaktuDibaca)
	if err != nil {

		return err
	}

	return nil

}
func (t *TekananImpln) GetAll(ctx context.Context) ([]domain.Data, error) {

	query := `
        SELECT id, device_id, tekanan_kpa, waktu_dibaca
        FROM tekanan.data
        ORDER BY waktu_dibaca DESC
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
			&data.TekananKpa,
			&data.WaktuDibaca,
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
        SELECT id, device_id, tekanan_kpa, waktu_dibaca
        FROM tekanan.data
        ORDER BY waktu_dibaca DESC
        LIMIT 1
    `

	var data domain.Data

	err := t.db.QueryRowContext(ctx, query).Scan(
		&data.ID,
		&data.DeviceId,
		&data.TekananKpa,
		&data.WaktuDibaca,
	)

	if err != nil {
		return domain.Data{}, err
	}

	return data, nil
}
