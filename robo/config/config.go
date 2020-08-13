package config

import (
	"flag"
	"sync"
)

var (
	load   sync.Once
	config = Config{}
)

// Config containing backend settings.
type Config struct {
	// MQTTBrokerURL points to a running Mosquitto server (our MQTT broker).
	// This is also where robots connect to publish their data.
	MQTTBrokerURL string `json:"mqtt_broker_url"`
	// MQTT topic that robots use to signal the start of a new
	// cleaning session.
	TopicRobotSessionStart string `json:"topic_robot_session_start"`
	// MQTT topic that robots use to signal the end of a cleaning
	// session.
	TopicRobotSessionEnd string `json:"topic_robot_session_end"`
	// MQTT topic that robots use to update their position during
	// a cleaning session.
	TopicRobotSessionUpdate string `json:"topic_robot_session_update"`

	// DgraphURL points to a running Dgraph server.
	DgraphURL string `json:"dgraph_url"`

	// DropAll flags that we want to drop the Dgraph schema
	// and recreate it.
	DropAll bool `json:"drop_all"`
	// Migrate flags that we want to apply Dgraph schema changes.
	Migrate bool `json:"migrate"`

	// APIURL is a URL pointing to a running instance of the backend, e.g.
	// https://api.example.com:10000
	APIURL string `json:"api_url"`
	// Host is the API hostname and port, e.g. api.example.com:3000
	Host string `json:"host"`
}

// GetConfig returns a singleton instance of the backend config.
func GetConfig() Config {
	load.Do(func() {
		flag.StringVar(&config.MQTTBrokerURL, "mqtt-broker-url", "tcp://localhost:1883", "set MQTT broker URL, e.g tcp://localhost:1883")

		flag.StringVar(&config.TopicRobotSessionStart, "topic-start", "/robot/session/start", "set MQTT topic for cleaning session start")
		flag.StringVar(&config.TopicRobotSessionEnd, "topic-end", "/robot/session/end", "set MQTT topic for cleaning session end")
		flag.StringVar(&config.TopicRobotSessionUpdate, "topic-update", "/robot/session/update", "set MQTT topic for robot session update")

		flag.BoolVar(&config.DropAll, "drop-all", false, "drop all tables and recreate schema")
		flag.BoolVar(&config.Migrate, "migrate", false, "migrate schema changes")

		config.DgraphURL = "127.0.0.1:9080"
		config.APIURL = "http://localhost:3000"
		config.Host = "localhost:3000"

		flag.Parse()
	})
	return config
}
