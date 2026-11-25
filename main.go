package main

import (
	"log"

	"menghitung-tekanan/internal/config"
	"menghitung-tekanan/internal/domain"
	tekananhandler "menghitung-tekanan/internal/handler/tekanan_handler"
	websockethandler "menghitung-tekanan/internal/handler/websocket"
	mqttsubscriber "menghitung-tekanan/internal/mqtt"
	"menghitung-tekanan/internal/repository"
	"menghitung-tekanan/internal/router"
	"menghitung-tekanan/internal/service"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/static"
	"github.com/joho/godotenv"
)

func main() {
	const mqttTopic = "bme280/data"
	// 1. Membaca konfigurasi dari .env.
	if err := godotenv.Load(); err != nil {
		log.Println("File .env tidak ditemukan, menggunakan environment variables")
	}

	cfg := config.Load()

	// 2. Menghubungkan PostgreSQL.
	db := config.Connect(cfg.Postgres)
	if db == nil {
		log.Println("Koneksi PostgreSQL gagal")
		return
	}
	defer db.Close()

	// 3. Menghubungkan MQTT Broker.
	mqttClient := config.InitMqtt()

	if mqttClient == nil || !mqttClient.IsConnected() {
		log.Println("Koneksi MQTT Broker gagal")
		return
	}
	defer mqttClient.Disconnect(250)

	// 4. Membuat Repository.
	tekananRepo := repository.NewTekananImpln(db)

	// 5. Membuat Service.
	tekananService := service.NewTekananService(tekananRepo)

	// 6. Membuat channel untuk data realtime.
	dataStream := make(chan domain.Data, 100)

	// 7. Membuat MQTT Subscriber.
	subscriber := mqttsubscriber.NewSubscriber(
		mqttClient,
		tekananService,
		dataStream,
	)

	// 8. Mulai subscribe MQTT.
	if err := subscriber.Subscribe(mqttTopic); err != nil {
		log.Printf("Gagal subscribe MQTT: %v", err)
		return
	}

	// 9. Membuat WebSocket Handler.
	wsHandler := websockethandler.NewTekananWebSocketHandler(
		dataStream,
	)

	// 10. Membuat HTTP Handler.
	httpHandler := tekananhandler.NewTekananHandler(
		tekananService,
	)

	// 11. Membuat Fiber App.
	app := fiber.New()
	app.Get("/", static.New("./public"))
	// 12. Mendaftarkan seluruh route.
	router.InitRouter(
		app,
		wsHandler,
		httpHandler,
	)

	// 13. Menjalankan server.
	log.Println("Server berjalan di http://localhost:3000")

	if err := app.Listen(":3000"); err != nil {
		log.Printf("Fiber berhenti: %v", err)
	}
}
