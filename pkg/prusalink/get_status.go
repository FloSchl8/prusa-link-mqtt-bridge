package prusalink

import (
	"encoding/json"
	"fmt"
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

func (c *Client) GetStatus() (*PrinterStatus, error) {
	body, err := c.Get("status")
	if err != nil {
		return nil, err
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
