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
	var req dto.DataReq

	if err := json.Unmarshal(msg.Payload(), &req); err != nil {
		s.log.WithFields(logrus.Fields{
			"operation": "decode MQTT",
			"layer":     "mqtt",
			"topic":     msg.Topic(),
		}).WithError(err).Error("payload MQTT tidak valid")
		return
	}

	deviceID := strings.TrimSpace(req.DeviceId)
	if deviceID == "" {
		deviceID = "bme280"
	}

	waktu := time.Now().UTC()
	if req.DataDibuat != nil && !req.DataDibuat.IsZero() {
		waktu = req.DataDibuat.UTC()
	}

	data := domain.Data{
		DeviceId:     deviceID,
		Temperatur:   req.Temperatur,
		Kelembapan:   req.Kelembapan,
		TekananUdara: req.TekananUdara,
		DataDibuat:   waktu,
	}

	if err := s.service.Save(context.Background(), data); err != nil {
		s.log.WithFields(logrus.Fields{
			"operation": "save data sensor",
			"layer":     "mqtt",
		}).WithError(err).Error("gagal menyimpan data BME280")
		return
	}

	s.log.WithFields(logrus.Fields{
		"temperatur":    data.Temperatur,
		"kelembapan":    data.Kelembapan,
		"tekanan_udara": data.TekananUdara,
	}).Info("data BME280 tersimpan dan diteruskan ke WebSocket")

	s.dataStream <- data
}
