package config

import (
	"menghitung-tekanan/internal/logger"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/sirupsen/logrus"
)

func InitMqtt() mqtt.Client {

	log := logger.New()
	opt := mqtt.NewClientOptions()
	opt.AddBroker("tcp://localhost:1883")
	client := mqtt.NewClient(opt)

	token := client.Connect()
	token.Wait()
	if token.Error() != nil {
		log.WithFields(logrus.Fields{
			"operation": "Init Mqtt",
			"layer":     "config",
		})
	}
	return client
}
