package prusalink

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetInfo(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/api/v1/info", r.URL.Path)
		assert.Equal(t, "test-api-key", r.Header.Get("X-Api-Key"))
		w.WriteHeader(http.StatusOK)
		fmt.Fprintln(w, `{"serial_number":"123456789"}`)
	}))
	defer server.Close()

	client := NewClient(server.URL[7:], "test-api-key")
	info, err := client.GetInfo()

	assert.NoError(t, err)
	assert.NotNil(t, info)
	assert.Equal(t, "123456789", info.SerialNumber)
}

func TestGetInfo_Error(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer server.Close()

	client := NewClient(server.URL[7:], "test-api-key")
	info, err := client.GetInfo()

	assert.Error(t, err)
	assert.Nil(t, info)
}
