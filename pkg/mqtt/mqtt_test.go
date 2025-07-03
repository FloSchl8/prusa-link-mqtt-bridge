package mqtt

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewMqttClient(t *testing.T) {
	// This is a basic test to ensure the client can be created.
	// A full integration test would require a running MQTT broker.
	_, err := NewMqttClient("localhost", "user", "pass", 1883)
	assert.Error(t, err) // Expect an error because no broker is running
}
