package entity

import (
	mqtt "github.com/eclipse/paho.mqtt.golang"
)

// MessageDelegator receives MQTT messages, converts
// them and passes them on to some other service.
type MessageDelegator interface {
	HandleStartSession(mqtt.Client, mqtt.Message)
	HandleUpdateSession(mqtt.Client, mqtt.Message)
	HandleEndSession(mqtt.Client, mqtt.Message)
}
