package router

import (
	tekananhandler "menghitung-tekanan/internal/handler/tekanan_handler"
	handler "menghitung-tekanan/internal/handler/websocket"
	"menghitung-tekanan/internal/middleware"

	"github.com/gofiber/contrib/v3/websocket"
	"github.com/gofiber/fiber/v3"
)

func InitRouter(app *fiber.App, t *handler.TekananWebSocketHandler, h *tekananhandler.TekananHandler) {

	app.Use("/ws", middleware.Update)
	app.Get("/ws/tekanan", websocket.New(t.Handle))

	//tampilkan html
	// HTTP
	app.Get("/api/tekanan/terbaru", h.AmbilTerbaru)
	app.Get("/api/tekanan/riwayat", h.AmbilRiwayat)
}
