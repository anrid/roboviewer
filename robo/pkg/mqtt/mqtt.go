package mqtt

import (
	"log"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

// Client is a thin wrapper around a standard MQTT client
// which we use to subscribe to messages being sent from
// robots.
type Client struct {
	c mqtt.Client
}

// NewClient connects to a MQTT broker and returns a Client
// instance.
func NewClient(broker string) *Client {
	if broker == "" {
		broker = "tcp://localhost:1883"
	}

	opts := mqtt.NewClientOptions().
		AddBroker(broker)
	client := mqtt.NewClient(opts)
	token := client.Connect()
	if token.Wait() && token.Error() != nil {
		log.Fatalf("could not connect to MQTT broker %s: %s", broker, token.Error())
	}
	return &Client{client}
}

// Subscribe to the given topic. Each new message on this topic will
// run the given handler.
func (c *Client) Subscribe(topic string, handler mqtt.MessageHandler) {
	token := c.c.Subscribe(topic, 0, handler)
	if token.Wait() && token.Error() != nil {
		log.Fatalf("could not subscribe to topic %s: %s", topic, token.Error())
	}
}

// Publish a message to the given topic.
func (c *Client) Publish(topic, message string) {
	token := c.c.Publish(topic, 0, false, message)
	if token.Wait() && token.Error() != nil {
		log.Fatalf("could not publish to topic %s: %s", topic, token.Error())
	}
}
