package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/disintegration/imaging"
)

type ImageResizeResponse struct {
	Status   string `json:"status"`
	Message  string `json:"message"`
	ImageURL string `json:"image_url"`
}

// MaxRequestSize is the maximum request size in bytes.
const (
	MaxRequestSize = 32 << 20 // 32 MB
	MaxImageSize   = 3500 * 3500
)

// isFormatSupported checks if a given image format is supported
func isFormatSupported(format string) bool {
	switch format {
	case "jpeg", "jpg", "png", "gif", "bmp", "tiff", "webp":
		return true
	default:
		return false
	}
}

// ImageResizeHandler is a handler for the /image-resize endpoint.
func ImageResizeHandler(w http.ResponseWriter, r *http.Request) {
	// Your API logic goes here// Create a new instance of ResponseWriterWrapper
	// and pass the response writer to it.
	rw := &ResponseWriterWrapper{
		w: w,
	}

	// Set the response content type to JSON
	rw.Header().Set("Content-Type", "application/json")

	// Check if it's a POST request
	if r.Method != http.MethodPost {
		// Set the response status code to 400 Bad Request
		rw.WriteHeader(http.StatusBadRequest)
		// Create and populate the response object
		response := ImageResizeResponse{
			Status:   "error",
			Message:  "Invalid request method",
			ImageURL: "",
		}
		// Encode the response object as JSON and write it to the response
		json.NewEncoder(rw).Encode(response)
		return
	}

	// Check if request size is too large
	if r.ContentLength > MaxRequestSize {
		// Set the response status code to 400 Bad Request
		rw.WriteHeader(http.StatusBadRequest)
		// Create and populate the response object
		response := ImageResizeResponse{
			Status:   "error",
			Message:  "Request size is too large",
			ImageURL: "",
		}
		// Encode the response object as JSON and write it to the response
		json.NewEncoder(rw).Encode(response)
		return
	}

	// Parse the form data
	err := r.ParseMultipartForm(MaxRequestSize) // Max 32 MB file size
	if err != nil {
		// Set the response status code to 400 Bad Request
		rw.WriteHeader(http.StatusBadRequest)
		// Create and populate the response object
		response := ImageResizeResponse{
			Status:   "error",
			Message:  "Failed to parse form data: " + err.Error(),
			ImageURL: "",
		}
		// Encode the response object as JSON and write it to the response
		json.NewEncoder(rw).Encode(response)
		return
	}

	// Get the format value from the form data
	format := r.FormValue("format")
	if format == "" {
		// Set the response status code to 400 Bad Request
		rw.WriteHeader(http.StatusBadRequest)
		// Create and populate the response object
		response := ImageResizeResponse{
			Status:   "error",
			Message:  "Missing format value",
			ImageURL: "",
		}
		// Encode the response object as JSON and write it to the response
		json.NewEncoder(rw).Encode(response)
		return
	}

	// Check if the form data "format" value is supported
	if !isFormatSupported(format) {
		// Set the response status code to 400 Bad Request
		rw.WriteHeader(http.StatusBadRequest)
		// Create and populate the response object
		response := ImageResizeResponse{
			Status:   "error",
			Message:  "Unsupported image format",
			ImageURL: "",
		}
		// Encode the response object as JSON and write it to the response
		json.NewEncoder(rw).Encode(response)
		return
	}

	// Get the uploaded file
	file, handler, err := r.FormFile("image")
	if err != nil {
		// Set the response status code to 400 Bad Request
		rw.WriteHeader(http.StatusBadRequest)
		// Create and populate the response object
		response := ImageResizeResponse{
			Status:   "error",
			Message:  "Failed to get uploaded file: " + err.Error(),
			ImageURL: "",
		}
		// Encode the response object as JSON and write it to the response
		json.NewEncoder(rw).Encode(response)
		return
	}
	defer file.Close()

	// Create a buffer to store the file data
	buf := make([]byte, handler.Size)

	// Read the file data into the buffer
	_, err = file.Read(buf)
	if err != nil {
		// Set the response status code to 400 Bad Request
		rw.WriteHeader(http.StatusBadRequest)
		// Create and populate the response object
		response := ImageResizeResponse{
			Status:   "error",
			Message:  "Failed to read file into buffer: " + err.Error(),
			ImageURL: "",
		}
		// Encode the response object as JSON and write it to the response
		json.NewEncoder(rw).Encode(response)
		return
	}

	// Create and populate the response object
	response := map[string]string{"message": "Hello, you reached the Image Resize handler!"}

	// Set the response content type to JSON
	w.Header().Set("Content-Type", "application/json")

	// Encode the response object as JSON and write it to the response
	// and return an error if the encoding fails.
	if err := json.NewEncoder(w).Encode(response); err != nil {
		// Set the response content type
		w.Header().Set("Content-Type", "text/plain")
		// Set the response status code to 500 Internal Server Error
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}


// resizeImage converts an image to a different format
func resizeImage(buf []byte, height string, width string) ([]byte, error) {
	// Use disintegration/imaging to resize the image type
	// and return the resized image as a byte slice.

	// Load the image from byte slice
	src, err := imaging.Decode(bytes.NewReader(buf))
	if err != nil {
		return nil, fmt.Errorf("failed to decode image: %w", err)
	}

	// Create a buffer to store the converted image
	outputBuf := bytes.NewBuffer(nil)

	// Resize the image with the provided dimensions
	// using the imaging package
	// TODO: Add support for custom dimensions

	// Return the resized image as a byte slice
	return outputBuf.Bytes(), nil
}