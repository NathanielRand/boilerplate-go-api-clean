package server

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/NathanielRand/webchest-image-converter-api/internal/routes"
)

func TestStartServer(t *testing.T) {
	// Get the router from the routes package
	router := routes.SetupRouter()

	// Create a test server that wraps the router
	testServer := httptest.NewServer(router)

	// Make a GET request to the /api/hello endpoint
	resp, err := http.Get(testServer.URL + "/api/v1/hello")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	defer resp.Body.Close()

	// Check that the response status code is 200 OK
	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status code %d, but got %d", http.StatusOK, resp.StatusCode)
	}
}
