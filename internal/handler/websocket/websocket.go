package handler

import (
	"log"

	"menghitung-tekanan/internal/domain"

	"github.com/gofiber/contrib/v3/websocket"
)

type wsClient struct {
	send chan domain.Data
}

type TekananWebSocketHandler struct {
	dataStream <-chan domain.Data

	register   chan *wsClient
	unregister chan *wsClient

	clients map[*wsClient]bool
}

func NewTekananWebSocketHandler(
	dataStream <-chan domain.Data,
) *TekananWebSocketHandler {
	h := &TekananWebSocketHandler{
		dataStream: dataStream,
		register:   make(chan *wsClient),
		unregister: make(chan *wsClient),
		clients:    make(map[*wsClient]bool),
	}

	// Menjalankan WebSocket Hub.
	go h.run()

	return h
}

// Hub menerima data dan membagikannya ke seluruh browser.
func (h *TekananWebSocketHandler) run() {
	for {
		select {

		case client := <-h.register:
			h.clients[client] = true

			log.Printf(
				"[WS] Browser terhubung. Total: %d",
				len(h.clients),
			)

		case client := <-h.unregister:
			if h.clients[client] {
				delete(h.clients, client)
				close(client.send)

				log.Printf(
					"[WS] Browser terputus. Total: %d",
					len(h.clients),
				)
			}

		case data, ok := <-h.dataStream:
			if !ok {
				log.Println("[WS] DataStream ditutup")
				return
			}

			// Broadcast ke semua browser.
			for client := range h.clients {
				select {
				case client.send <- data:
					// Data berhasil masuk antrean browser.

				default:
					log.Println(
						"[WS] Antrean salah satu browser penuh",
					)
				}
			}
		}
	}
}

// Menangani satu koneksi WebSocket.
func (h *TekananWebSocketHandler) Handle(conn *websocket.Conn) {
	client := &wsClient{
		send: make(chan domain.Data, 100),
	}

	// Daftarkan browser ke Hub.
	h.register <- client

	defer func() {
		h.unregister <- client
		conn.Close()
	}()

	clientDone := make(chan struct{})

	// Mendeteksi jika browser terputus.
	go func() {
		defer close(clientDone)

		for {
			if _, _, err := conn.ReadMessage(); err != nil {
				return
			}
		}
	}()

	// Mengirim data dari antrean khusus browser ini.
	for {
		select {
		case <-clientDone:
			return

		case data, ok := <-client.send:
			if !ok {
				return
			}

			if err := conn.WriteJSON(data); err != nil {
				log.Println("[WS] Gagal mengirim data:", err)
				return
			}
		}
	}
}
