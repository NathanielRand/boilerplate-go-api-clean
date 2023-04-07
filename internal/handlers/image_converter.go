package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

// ImageConvertHandler is a handler for the /convert/image/ endpoint.
// It converts an image to a different format. It accepts a POST request
// with a form data file named "image" and a URL path with the from/to format values.
func ImageConvertHandler(w http.ResponseWriter, r *http.Request) {
	// Check if it's a POST request
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, `{"status": "error", "message": "Invalid request method"}`)
		return
	}

	// Extract from and to values from URL slug
	from, to := extractFromToValuesFromURL(r)
	print("FROM: ", from)
	print("TO: ", to)

	// Check if the URL from/to format value is supported
	if !isFormatSupported(from) || !isFormatSupported(to) {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, `{"status": "error", "message": "Unsupported image format"}`)
		return
	}

	// Parse the form data
	err := r.ParseMultipartForm(32 << 20) // Max 32 MB file size
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, `{"status": "error", "message": "Failed to parse form data. File size may be too large or the form data may be invalid", "error": "%s"}`, err)
		return
	}

	// Get the uploaded file
	file, handler, err := r.FormFile("image")
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, `{"status": "error", "message": "Failed to get uploaded file", "error": "%s"}`, err)
		return
	}
	defer file.Close()

	// Create a buffer to store the file data
	buf := make([]byte, handler.Size)

	// Read the file data into the buffer
	_, err = file.Read(buf)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, `{"status": "error", "message": "Failed to read file into buffer", "error": "%s"}`, err)
		return
	}

	// Convert the image
	convertedFile, err := convertImage(buf, from, to)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, `{"status": "error", "message": "Failed to convert image", "error": "%s"}`, err)
		return
	}

	// Store the converted image in cloud storage
	// and return an authenicated URL to the converted image.
	convertedImageURL, err := storeImage(convertedFile)

	// Create and populate the response object
	response := map[string]string{
		"status": "success",
		"status_code": "200",
		"image_url": convertedImageURL,
		"message": "Image converted successfully",
	}

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

// convertImage converts an image to a different format
func convertImage(buf []byte, from, to string) ([]byte, error) {
	// TODO: Use t
}

// storeImage stores an image in cloud storage and 
// returns an authenticated URL to the image
func storeImage(buf []byte) (string, error) {
	// TODO: Call cloud storage service to store the image
	// and return an authenticated URL to the image
	
}

// Function to extract from and to values from URL slug
func extractFromToValuesFromURL(r *http.Request) (from, to string) {
	// Extract from and to values from URL slug
	slug := strings.TrimPrefix(r.URL.Path, "/api/v1/convert/image/")
	// Split the slug using '/' as delimiter
	parts := strings.Split(slug, "/")
	// Extract from and to values from the parts slice
	from = parts[0]
	to = parts[1]
	// Return the from and to values
	return from, to
}

// isFormatSupported checks if a given image format is supported
func isFormatSupported(format string) bool {
	// Supported image formats
	supportedFormats := []string{"any", "jpg", "jpeg", "png", "bmp", "gif", "tiff", "webp"}
	// Check if the format is supported
	for _, f := range supportedFormats {
		if f == format {
			return true
		}
	}
	return false
}
