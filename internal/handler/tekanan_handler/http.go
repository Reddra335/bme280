package tekananhandler

import (
	"context"
	"database/sql"
	"errors"
	"menghitung-tekanan/internal/domain"

	"github.com/gofiber/fiber/v3"
)

type TekananService interface {
	Save(ctx context.Context, data domain.Data) error
	GetLatest(ctx context.Context) (domain.Data, error)
	GetAll(ctx context.Context) ([]domain.Data, error)
}

type TekananHandler struct {
	service TekananService
}

func NewTekananHandler(service TekananService) *TekananHandler {
	return &TekananHandler{
		service: service,
	}
}

// Mengambil tekanan terbaru
func (h *TekananHandler) AmbilTerbaru(c fiber.Ctx) error {
	data, err := h.service.GetLatest(c.Context())

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"pesan": "Belum ada data tekanan",
			})
		}

		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"pesan": "Gagal mengambil data tekanan",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"pesan": "Berhasil mengambil tekanan terbaru",
		"data":  data,
	})
}

// Mengambil seluruh history tekanan
func (h *TekananHandler) AmbilRiwayat(c fiber.Ctx) error {
	data, err := h.service.GetAll(c.Context())

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"pesan": "Gagal mengambil history tekanan",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"pesan": "Berhasil mengambil history tekanan",
		"data":  data,
	})
}
