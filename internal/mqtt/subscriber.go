package mqtt

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"menghitung-tekanan/internal/domain"
	"menghitung-tekanan/internal/dto"
	"menghitung-tekanan/internal/logger"

	paho "github.com/eclipse/paho.mqtt.golang"
	"github.com/sirupsen/logrus"
)

// Service yang dibutuhkan subscriber.
type TekananService interface {
	Save(ctx context.Context, data domain.Data) error
}

type Subscriber struct {
	client     paho.Client
	service    TekananService
	dataStream chan<- domain.Data
	log        *logger.Logger
}

func NewSubscriber(
	client paho.Client,
	service TekananService,
	dataStream chan<- domain.Data,
) *Subscriber {
	return &Subscriber{
		client:     client,
		service:    service,
		dataStream: dataStream,
		log:        logger.New(),
	}
}

// Berlangganan ke topic MQTT.
func (s *Subscriber) Subscribe(topic string) error {
	if strings.TrimSpace(topic) == "" {
		return errors.New("topic MQTT tidak boleh kosong")
	}

	token := s.client.Subscribe(topic, 0, s.handleMessage)
	token.Wait()

	if err := token.Error(); err != nil {
		return fmt.Errorf("gagal subscribe MQTT: %w", err)
	}

	s.log.Info("berhasil subscribe MQTT: " + topic)

	return nil
}

// Menerima pesan dari MQTT Broker.
func (s *Subscriber) handleMessage(client paho.Client, msg paho.Message) {
	// 1. Decode JSON menjadi DTO.
	var req dto.DataReq

	if err := json.Unmarshal(msg.Payload(), &req); err != nil {
		s.log.WithFields(logrus.Fields{
			"operation": "decode MQTT",
			"layer":     "mqtt",
			"topic":     msg.Topic(),
		}).WithError(err).Error("payload MQTT tidak valid")

		return
	}

	// 2. Validasi device ID.
	deviceID := strings.TrimSpace(req.DeviceId)

	if deviceID == "" {
		s.log.Error("device_id tidak boleh kosong")
		return
	}

	// 3. Ubah DTO menjadi domain.
	data := domain.Data{
		DeviceId:    deviceID,
		TekananKpa:  req.TekananKpa,
		WaktuDibaca: time.Now().UTC(),
	}

	s.log.WithFields(logrus.Fields{
		"topic":       msg.Topic(),
		"device_id":   data.DeviceId,
		"tekanan_kpa": data.TekananKpa,
	}).Info("pesan MQTT berhasil diterima")

	// 4. Simpan ke PostgreSQL melalui service.
	if err := s.service.Save(context.Background(), data); err != nil {
		s.log.WithFields(logrus.Fields{
			"operation": "save tekanan",
			"layer":     "mqtt",
		}).WithError(err).Error("gagal menyimpan data tekanan")

		return
	}

	s.log.Info(fmt.Sprintf(
		"data berhasil disimpan ke PostgreSQL: %.2f kPa",
		data.TekananKpa,
	))

	// 5. Teruskan data yang sudah berhasil disimpan.
	// Tidak menggunakan default agar data tidak langsung dibuang.
	s.dataStream <- data

	s.log.Info(fmt.Sprintf(
		"data berhasil diteruskan ke channel WebSocket: %.2f kPa",
		data.TekananKpa,
	))
}
