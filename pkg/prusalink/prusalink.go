package prusalink

import (
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"time"
)

// APIResponse matches the nested structure of the PrusaLink API spec.

type APIResponse struct {
	Job     *JobStatus   `json:"job,omitempty"`
	Printer *PrinterInfo `json:"printer"`
}

type JobStatus struct {
	Progress      float64 `json:"progress"`
	TimeRemaining int     `json:"time_remaining"`
	TimePrinting  int     `json:"time_printing"`
}

type PrinterInfo struct {
	State        string  `json:"state"`
	TempBed      float64 `json:"temp_bed"`
	TargetBed    float64 `json:"target_bed"`
	TempNozzle   float64 `json:"temp_nozzle"`
	TargetNozzle float64 `json:"target_nozzle"`
	AxisZ        float64 `json:"axis_z"`
	Flow         int     `json:"flow"`
	Speed        int     `json:"speed"`
	FanHotend    int     `json:"fan_hotend"`
	FanPrint     int     `json:"fan_print"`
}

// PrinterStatus is the flattened structure we publish to MQTT.

type PrinterStatus struct {
	StateText     string  `json:"state_text"`
	TempNozzle    float64 `json:"temp_nozzle"`
	TargetNozzle  float64 `json:"target_nozzle"`
	TempBed       float64 `json:"temp_bed"`
	TargetBed     float64 `json:"target_bed"`
	AxisZ         float64 `json:"axis_z"`
	Flow          int     `json:"flow"`
	Speed         int     `json:"speed"`
	FanHotend     int     `json:"fan_hotend"`
	FanPrint      int     `json:"fan_print"`
	Progress      float64 `json:"progress"`
	PrintTimeLeft int     `json:"print_time_left"`
	PrintTime     int     `json:"print_time"`
}

func GetStatus(host, apiKey string) (*PrinterStatus, error) {
	client := &http.Client{
		Timeout: 5 * time.Second,
	}
	slog.Debug("Attempting to get status from PrusaLink", "host", host)
	req, err := http.NewRequest("GET", fmt.Sprintf("http://%s/api/v1/status", host), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("X-Api-Key", apiKey)

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	slog.Debug("Received response from PrusaLink API",
		"status_code", resp.StatusCode,
		"body", string(body),
	)

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var apiResponse APIResponse
	if err := json.Unmarshal(body, &apiResponse); err != nil {
		return nil, fmt.Errorf("failed to unmarshal API response: %w", err)
	}

	// Map the nested API structure to our flat MQTT structure.
	status := &PrinterStatus{}
	if apiResponse.Printer != nil {
		status.StateText = apiResponse.Printer.State
		status.TempNozzle = apiResponse.Printer.TempNozzle
		status.TargetNozzle = apiResponse.Printer.TargetNozzle
		status.TempBed = apiResponse.Printer.TempBed
		status.TargetBed = apiResponse.Printer.TargetBed
		status.AxisZ = apiResponse.Printer.AxisZ
		status.Flow = apiResponse.Printer.Flow
		status.Speed = apiResponse.Printer.Speed
		status.FanHotend = apiResponse.Printer.FanHotend
		status.FanPrint = apiResponse.Printer.FanPrint
	}
	if apiResponse.Job != nil {
		status.Progress = apiResponse.Job.Progress
		status.PrintTimeLeft = apiResponse.Job.TimeRemaining
		status.PrintTime = apiResponse.Job.TimePrinting
	}

	return status, nil
}
