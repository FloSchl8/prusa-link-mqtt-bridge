package main

import (
	"context"
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

	mqttClient, err := mqtt.NewMqttClient(cfg.Mqtt.Broker, cfg.Mqtt.Username, cfg.Mqtt.Password, cfg.Mqtt.Port)
	if err != nil {
		slog.Error("Failed to connect to MQTT broker", "error", err)
		os.Exit(1)
	}
	slog.Info("Connected to MQTT broker")

	ticker := time.NewTicker(time.Duration(cfg.PrusaLink.Interval) * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		slog.Info("Ticker ticked, fetching status...")
		status, err := prusalink.GetStatus(cfg.PrusaLink.Host, cfg.PrusaLink.ApiKey)
		if err != nil {
			slog.Error("Failed to get printer status", "error", err)
			continue
		}

		if err := mqtt.PublishStatus(mqttClient, cfg.Mqtt.Topic, status); err != nil {
			slog.Error("Failed to publish status to MQTT", "error", err)
		}
	}
}
