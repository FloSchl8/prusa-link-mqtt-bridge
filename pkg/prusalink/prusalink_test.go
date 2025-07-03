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
			fmt.Fprintln(w, `{"state_text": "Printing", "temp_nozzle": 215, "target_nozzle": 215, "temp_bed": 85, "target_bed": 85, "progress": 50.5}`)
		}))
		defer server.Close()

		status, err := GetStatus(server.URL[7:], "test-api-key")
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

		_, err := GetStatus(server.URL[7:], "wrong-api-key")
		assert.Error(t, err)
	})
}
