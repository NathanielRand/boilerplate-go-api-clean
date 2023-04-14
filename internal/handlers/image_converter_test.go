package handlers_test

import (
	"bytes"
	"encoding/json"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/NathanielRand/webchest-image-converter-api/internal/handlers"
)

func TestImageConvertHandler(t *testing.T) {
	// Create a mock HTTP request with a POST method, form data, and file
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	file, err := writer.CreateFormFile("image", "test.jpg")
	if err != nil {
		t.Fatalf("Failed to create form file: %v", err)
	}
	// Add test image file data to the request
	file.Write([]byte("test image file data"))
	writer.Close()
	req := httptest.NewRequest("POST", "/convert/image/jpeg", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())

	// Create a mock HTTP response recorder
	recorder := httptest.NewRecorder()

	// Call the ImageConvertHandler with the mock request and response recorder
	handlers.ImageConvertHandler(recorder, req)

	// Check the response status code
	if recorder.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, recorder.Code)
	}

	// Check the response body
	var response handlers.ImageConvertResponse
	err = json.NewDecoder(recorder.Body).Decode(&response)
	if err != nil {
		t.Fatalf("Failed to decode response body: %v", err)
	}

	// Check the response fields
	if response.Status != "success" {
		t.Errorf("Expected response status 'success', got '%s'", response.Status)
	}
	if response.Message != "" {
		t.Errorf("Expected empty response message, got '%s'", response.Message)
	}
	if response.ImageURL == "" {
		t.Error("Expected non-empty image URL")
	}
}
