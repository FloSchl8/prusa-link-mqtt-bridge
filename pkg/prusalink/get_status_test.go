package prusalink

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetStatus(t *testing.T) {
	t.Run("successful request", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "/api/v1/status", r.URL.Path)
			assert.Equal(t, "test-api-key", r.Header.Get("X-Api-Key"))
			w.WriteHeader(http.StatusOK)
			fmt.Fprintln(w, `{"printer":{"state":"Printing","temp_nozzle":215,"target_nozzle":215,"temp_bed":85,"target_bed":85,"axis_z":10.5,"flow":100,"speed":100,"fan_hotend":255,"fan_print":255},"job":{"progress":50.5,"time_remaining":120,"time_printing":60}}`)
		}))
		defer server.Close()

		client := NewClient(server.URL[7:], "test-api-key")
		status, err := client.GetStatus()
		assert.NoError(t, err)
		assert.NotNil(t, status)
		assert.Equal(t, "Printing", status.StateText)
		assert.Equal(t, 215.0, status.TempNozzle)
		assert.Equal(t, 50.5, status.Progress)
	})

	t.Run("API error", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusUnauthorized)
		}))
		defer server.Close()

		client := NewClient(server.URL[7:], "wrong-api-key")
		_, err := client.GetStatus()
		assert.Error(t, err)
	})

	t.Run("server error", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
		}))
		defer server.Close()

		client := NewClient(server.URL[7:], "test-api-key")
		_, err := client.GetStatus()
		assert.Error(t, err)
	})
}
