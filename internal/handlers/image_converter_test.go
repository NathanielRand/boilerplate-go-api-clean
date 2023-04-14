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

// BenchmarkImageConvertHandler benchmarks the ImageConvertHandler function
func BenchmarkImageConvertHandler(b *testing.B) {
	// Create a mock HTTP request with a POST method, form data, and file
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	format, err := writer.CreateFormField("format")
	if err != nil {
		b.Fatalf("Failed to create form field: %v", err)
	}
	if _, err = format.Write([]byte("png")); err != nil {
		b.Fatalf("Failed to write to form field: %v", err)
	}
	file, err := writer.CreateFormFile("image", "test.jpg")
	if err != nil {
		b.Fatalf("Failed to create form file: %v", err)
	}
	// Add test image file data to the request
	file.Write([]byte("test image file data"))
	writer.Close()
	req := httptest.NewRequest("POST", "/convert/image", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.Header.Set("X-RapidAPI-Proxy-Secret", "78f5b3e0-d3d0-11ed-bf92-43930995aeef")

	// Create a mock HTTP response recorder
	recorder := httptest.NewRecorder()

	// Call the ImageConvertHandler with the mock request and response recorder
	for i := 0; i < b.N; i++ {
		handlers.ImageConvertHandler(recorder, req)
	}
}

func TestImageConvertHandler(t *testing.T) {
	// Create a mock HTTP request with a POST method, form data, and file
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	format, err := writer.CreateFormField("format")
	if err != nil {
		t.Fatalf("Failed to create form field: %v", err)
	}
	if _, err = format.Write([]byte("png")); err != nil {
		t.Fatalf("Failed to write to form field: %v", err)
	}
	file, err := writer.CreateFormFile("image", "test.jpg")
	if err != nil {
		t.Fatalf("Failed to create form file: %v", err)
	}
	// Add test image file data to the request
	file.Write([]byte("test image file data"))
	writer.Close()
	req := httptest.NewRequest("POST", "/convert/image", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.Header.Set("X-RapidAPI-Proxy-Secret", "78f5b3e0-d3d0-11ed-bf92-43930995aeef")

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
