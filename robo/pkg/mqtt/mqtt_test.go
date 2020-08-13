package mqtt

import (
	"fmt"
	"sync"
	"testing"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

func TestMqttPubSub(t *testing.T) {
	const TOPIC = "mytopic/test"

	c := NewClient("tcp://localhost:1883")

	var wg sync.WaitGroup
	wg.Add(1)

	c.Subscribe(TOPIC, func(client mqtt.Client, msg mqtt.Message) {
		if string(msg.Payload()) != "mymessage" {
			t.Fatalf("want mymessage, got %s", msg.Payload())
		}
		fmt.Printf("message payload: %s\n", msg.Payload())

		wg.Done()
	})

	c.Publish(TOPIC, "mymessage")

	wg.Wait()
}
