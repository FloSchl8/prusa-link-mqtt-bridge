package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"strings"
	"time"

	"prusa-link-mqtt-bridge/pkg/config"
	"prusa-link-mqtt-bridge/pkg/mqtt"
	"prusa-link-mqtt-bridge/pkg/prusalink"
)

func main() {
	ctx := context.Background()
	cfg, err := config.LoadConfig(ctx)
	if err != nil {
		slog.Error("Failed to load config", "error", err)
		os.Exit(1)
	}

	// Setup logger
	var level slog.Level
	switch strings.ToUpper(cfg.Log.Level) {
	case "DEBUG":
		level = slog.LevelDebug
	case "INFO":
		level = slog.LevelInfo
	case "WARN":
		level = slog.LevelWarn
	case "ERROR":
		level = slog.LevelError
	default:
		slog.Warn("Invalid log level specified, defaulting to INFO", "level", cfg.Log.Level)
		level = slog.LevelInfo
	}

	var handler slog.Handler
	if cfg.Log.Format == "json" {
		handler = slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: level})
	} else {
		handler = slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: level})
	}
	slog.SetDefault(slog.New(handler))

	slog.Info("Configuration loaded successfully")

	prusaClient := prusalink.NewClient(cfg.PrusaLink.Host, cfg.PrusaLink.ApiKey)
	info, err := prusaClient.GetInfo()
	if err != nil {
		slog.Error("Failed to get printer info", "error", err)
		os.Exit(1)
	}
	slog.Info("Got printer info", "serial_number", info.SerialNumber)

	availabilityTopic := fmt.Sprintf("%s/%s/status", cfg.Mqtt.Topic, info.SerialNumber)
	statusTopic := fmt.Sprintf("%s/%s/status", cfg.Mqtt.Topic, info.SerialNumber)

	mqttClient, err := mqtt.NewMqttClient(cfg.Mqtt.Broker, cfg.Mqtt.Username, cfg.Mqtt.Password, cfg.Mqtt.Port, availabilityTopic)
	if err != nil {
		slog.Error("Failed to connect to MQTT broker", "error", err)
		os.Exit(1)
	}
	slog.Info("Connected to MQTT broker")

	if err := mqtt.Publish(mqttClient, availabilityTopic, "online"); err != nil {
		slog.Error("Failed to publish online status", "error", err)
	}

	ticker := time.NewTicker(time.Duration(cfg.PrusaLink.Interval) * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		slog.Info("Ticker ticked, fetching status...")
		status, err := prusaClient.GetStatus()
		if err != nil {
			slog.Error("Failed to get printer status", "error", err)
			if err := mqtt.Publish(mqttClient, availabilityTopic, "offline"); err != nil {
				slog.Error("Failed to publish offline status", "error", err)
			}
			continue
		}

		if err := mqtt.PublishStatus(mqttClient, statusTopic, status); err != nil {
			slog.Error("Failed to publish status to MQTT", "error", err)
		}
	}
}
