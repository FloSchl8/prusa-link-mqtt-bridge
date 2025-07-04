package prusalink

import "encoding/json"

// APIInfoResponse matches the structure of the /api/v1/info endpoint.
type APIInfoResponse struct {
	SerialNumber string `json:"serial_number"`
}

func (c *Client) GetInfo() (*APIInfoResponse, error) {
	body, err := c.Get("info")
	if err != nil {
		return nil, err
	}

	var info APIInfoResponse
	if err := json.Unmarshal(body, &info); err != nil {
		return nil, err
	}

	return &info, nil
}
