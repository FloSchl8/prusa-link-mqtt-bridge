package config

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoadConfig(t *testing.T) {
	t.Run("loads config from environment variables", func(t *testing.T) {
		t.Setenv("PRUSALINK_HOST", "testhost")
		t.Setenv("PRUSALINK_APIKEY", "testkey")
		t.Setenv("PRUSALINK_INTERVAL", "15")
		t.Setenv("MQTT_BROKER", "testbroker")
		t.Setenv("MQTT_PORT", "1884")
		t.Setenv("MQTT_USERNAME", "testuser")
		t.Setenv("MQTT_PASSWORD", "testpass")
		t.Setenv("MQTT_TOPIC", "test-topic")

		cfg, err := LoadConfig(context.Background())
		assert.NoError(t, err)

		assert.Equal(t, "testhost", cfg.PrusaLink.Host)
		assert.Equal(t, "testkey", cfg.PrusaLink.ApiKey)
		assert.Equal(t, 15, cfg.PrusaLink.Interval)
		assert.Equal(t, "testbroker", cfg.Mqtt.Broker)
		assert.Equal(t, 1884, cfg.Mqtt.Port)
		assert.Equal(t, "testuser", cfg.Mqtt.Username)
		assert.Equal(t, "testpass", cfg.Mqtt.Password)
		assert.Equal(t, "test-topic", cfg.Mqtt.Topic)
	})

	t.Run("returns error if required variables are missing", func(t *testing.T) {
		// Unset required variables
		t.Setenv("PRUSALINK_HOST", "")
		t.Setenv("PRUSALINK_APIKEY", "")
		t.Setenv("MQTT_BROKER", "")

		_, err := LoadConfig(context.Background())
		assert.Error(t, err)
	})

	t.Run("uses default values", func(t *testing.T) {
		t.Setenv("PRUSALINK_HOST", "testhost")
		t.Setenv("PRUSALINK_APIKEY", "testkey")
		t.Setenv("MQTT_BROKER", "testbroker")

		cfg, err := LoadConfig(context.Background())
		assert.NoError(t, err)

		assert.Equal(t, 10, cfg.PrusaLink.Interval)
		assert.Equal(t, 1883, cfg.Mqtt.Port)
		assert.Equal(t, "prusa-link", cfg.Mqtt.Topic)
	})
}
