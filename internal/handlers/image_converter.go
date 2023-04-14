package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"net/http"
	"time"

	// Internal
	"github.com/NathanielRand/webchest-image-converter-api/internal/config"
	"github.com/NathanielRand/webchest-image-converter-api/internal/repositories"

	// External
	"github.com/chai2010/webp"
	"github.com/disintegration/imaging"
)

// ImageConvertResponse is the response object for the /convert/image/ endpoint.
type ImageConvertResponse struct {
	Status   string `json:"status"`
	Message  string `json:"message"`
	ImageURL string `json:"image_url"`
}

// MaxRequestSize is the maximum request size in bytes.
const (
	MaxRequestSize = 32 << 20 // 32 MB
	MaxImageSize   = 3500 * 3500
)

// Error messages
var (
	ErrInvalidRequestMethod = errors.New("Invalid request method")
	ErrRequestSizeTooLarge  = errors.New("Request size is too large")
	ErrParseFormDataFailed  = errors.New("Failed to parse form data")
	ErrMissingFormatValue   = errors.New("Missing format value")
	ErrUnsupportedFormat    = errors.New("Unsupported image format")
	ErrFailedToGetFile      = errors.New("Failed to get uploaded file")
)

// formatMapping is a map of supported image formats.
// The key is the format name and the value is the imaging.Format value.
// See https://godoc.org/github.com/disintegration/imaging#Format
var formatMapping = map[string]imaging.Format{
	"jpeg": imaging.JPEG,
	"jpg":  imaging.JPEG,
	"png":  imaging.PNG,
	"gif":  imaging.GIF,
	"bmp":  imaging.BMP,
	"tiff": imaging.TIFF,
	// webp is not supported by imaging.Format
	// check implementation below for webp support.
}

// isFormatSupported checks if a given image format is supported
func isFormatSupported(format string) bool {
	switch format {
	case "jpeg", "jpg", "png", "gif", "bmp", "tiff", "webp":
		return true
	default:
		return false
	}
}

// ImageConvertHandler is a handler for the /convert/image/ endpoint.
// It converts an image to a different format. It accepts a POST request
// with a form data file named "image" and a URL path with the from/to format values.
func ImageConvertHandler(w http.ResponseWriter, r *http.Request) {
	// Create a new instance of ResponseWriterWrapper
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
		response := ImageConvertResponse{
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
		response := ImageConvertResponse{
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
		response := ImageConvertResponse{
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
		response := ImageConvertResponse{
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
		response := ImageConvertResponse{
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
		response := ImageConvertResponse{
			Status:   "error",
			Message:  "Failed to get uploaded file: " + err.Error(),
			ImageURL: "",
		}
		// Encode the response object as JSON and write it to the response
		json.NewEncoder(rw).Encode(response)
		return
	}
	defer file.Close()

	// Check image size, if over 3500x3500, reject the request
	// and return an error.
	// if handler.Size > 3500*3500 {
	// 	// Set the response status code to 400 Bad Request
	// 	w.WriteHeader(http.StatusBadRequest)
	// 	// Create and populate the response object
	// 	response := ImageConvertResponse{
	// 		Status:   "error",
	// 		Message:  "Image size is too large. Try using a smaller image or resize your image with our /resize/image endpoint",
	// 		ImageURL: "",
	// 	}
	// 	// Encode the response object as JSON and write it to the response
	// 	json.NewEncoder(w).Encode(response)
	// 	return
	// }

	// Create a buffer to store the file data
	buf := make([]byte, handler.Size)

	// Read the file data into the buffer
	_, err = file.Read(buf)
	if err != nil {
		// Set the response status code to 400 Bad Request
		rw.WriteHeader(http.StatusBadRequest)
		// Create and populate the response object
		response := ImageConvertResponse{
			Status:   "error",
			Message:  "Failed to read file into buffer: " + err.Error(),
			ImageURL: "",
		}
		// Encode the response object as JSON and write it to the response
		json.NewEncoder(rw).Encode(response)
		return
	}

	// Convert the image
	convertedFile, err := convertImage(buf, format)
	if err != nil {
		// Set the response status code to 500 Internal Server Error
		rw.WriteHeader(http.StatusInternalServerError)
		// Create and populate the response object
		response := ImageConvertResponse{
			Status:   "error",
			Message:  "Failed to convert image: " + err.Error(),
			ImageURL: "",
		}
		// Encode the response object as JSON and write it to the response
		json.NewEncoder(rw).Encode(response)
		return
	}

	// Store the converted image in cloud storage
	// and return an authenticated URL to the converted image.
	convertedImageURL, err := storeImage(convertedFile, format)
	if err != nil {
		// Set the response status code to 500 Internal Server Error
		rw.WriteHeader(http.StatusInternalServerError)
		// Create and populate the response object
		response := ImageConvertResponse{
			Status:   "error",
			Message:  "Failed to store image: " + err.Error(),
			ImageURL: "",
		}
		// Encode the response object as JSON and write it to the response
		json.NewEncoder(rw).Encode(response)
		return
	}

	// Create and populate the response object
	response := ImageConvertResponse{
		Status:   "success",
		Message:  "Image successfully converted to " + format + " format",
		ImageURL: convertedImageURL,
	}

	// Encode the response object as JSON and write it to the response
	// and return an error if the encoding fails.
	if err := json.NewEncoder(rw).Encode(response); err != nil {
		// Set the response status code to 500 Internal Server Error
		rw.WriteHeader(http.StatusInternalServerError)
		// Create and populate the response object
		response := ImageConvertResponse{
			Status:   "error",
			Message:  "Failed to encode response as JSON: " + err.Error(),
			ImageURL: "",
		}
		// Encode the response object as JSON and write it to the response
		json.NewEncoder(rw).Encode(response)
		return
	}
}

// convertImage converts an image to a different format
func convertImage(buf []byte, format string) ([]byte, error) {
	// Use disintegration/imaging to convert the image type
	// and return the converted image as a byte slice.

	// Load the image from byte slice
	src, err := imaging.Decode(bytes.NewReader(buf))
	if err != nil {
		return nil, fmt.Errorf("failed to decode image: %w", err)
	}

	// // Create a buffer to store the converted image
	outputBuf := bytes.NewBuffer(nil)

	// DEFINITELY NOT THE BEST WAY TO DO THIS. vvvvv
	// REFACTOR THIS LATER.

	// If the image is being converted to WEBP,
	// use chai2010/webp to convert the image to WEBP
	// and return the converted image as a byte slice.
	if format == "webp" {
		// Convert the image to WEBP
		// and return the converted image as a byte slice.
		// Encode lossless webp
		if err = webp.Encode(outputBuf, src, &webp.Options{Lossless: true}); err != nil {
			return nil, fmt.Errorf("failed to encode image to WEBP: %w", err)
		}
	} else {
		// Save the converted image to the buffer in the desired format
		err = imaging.Encode(outputBuf, src, formatMapping[format])
		if err != nil {
			return nil, fmt.Errorf("failed to encode image to %s: %w", format, err)
		}
	}

	return outputBuf.Bytes(), nil
}

// storeImage stores an image in cloud storage and returns an authenticated URL to the image
func storeImage(buf []byte, newExt string) (string, error) {
	// Generate a random URL with timestamp
	rand.Seed(time.Now().UnixNano())
	randNum := rand.Intn(1000)
	timestamp := time.Now().Format("20060102150405")
	name := fmt.Sprintf("webchest_%d_image_%s.%s", randNum, timestamp, newExt)

	// Get the context
	ctx := context.Background()

	// Create a new context with a timeout of 10 seconds
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel() // make sure to cancel the context to avoid potential resource leaks

	// Bucket name
	bucketName := "webchest_image_converter"

	// Get the global storage client instance
	client := config.GetStorageClient()

	// Create a new instance of CloudStorageRepository
	cloudStorageRepo, err := repositories.NewCloudStorageRepository(bucketName, client)
	if err != nil {
		return "", err
	}

	// Use a Goroutine to call UploadImage on the cloudStorageRepo instance to upload the image to cloud storage
	// and return an authenticated URL to the image.
	ch := make(chan string, 1) // channel to receive the result of the Goroutine
	go func() {
		defer close(ch) // close the channel when the Goroutine completes

		url, err := cloudStorageRepo.UploadImage(ctx, name, bytes.NewReader(buf))
		if err != nil {
			ch <- err.Error() // send the error to the channel
			return
		}
		ch <- url // send the URL to the channel
	}()

	// Return the channel as a future that will contain the result of the Goroutine
	select {
	case result := <-ch:
		return result, nil
	case <-ctx.Done():
		return "", ctx.Err() // return the context error if the Goroutine takes longer than the timeout
	}
}
