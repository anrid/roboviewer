package config

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestConfig(t *testing.T) {
	c := GetConfig()
	require.Contains(t, c.MQTTBrokerURL, "tcp://", "should contain the default value 'tcp://'")
	require.Contains(t, c.TopicRobotSessionStart, "/robot/session", "should contain the default value '/robot/session/...'")
}
