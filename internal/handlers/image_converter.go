package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	// "image"
	// "image/color"
	// "image/draw"
	"image/png"
	"log"
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

// formatMapping is a map of supported image formats.
// The key is the format name and the value is the imaging.Format value.
// See https://godoc.org/github.com/disintegration/imaging#Format
var formatMapping = map[string]imaging.Format{
	"jpeg": imaging.JPEG,
	"png":  imaging.PNG,
	"gif":  imaging.GIF,
	"bmp":  imaging.BMP,
	"tiff": imaging.TIFF,
	// "webp": imaging.WEBP,
}

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

	// Check if request size is too large
	if r.ContentLength > config.MaxRequestSize() {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, `{"status": "error", "message": "Request size is too large"}`)
		return
	}

	// Parse the form data
	err := r.ParseMultipartForm(32 << 20) // Max 32 MB file size
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, `{"status": "error", "message": "Failed to parse form data. File size may be too large or the form data may be invalid", "error": "%s"}`, err)
		return
	}

	// Get the format value from the form data
	format := r.FormValue("format")
	if format == "" {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, `{"status": "error", "message": "Missing format value"}`)
		return
	}

	// Check if the URL from/to format value is supported
	if !isFormatSupported(format) {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, `{"status": "error", "message": "Unsupported image format"}`)
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

	// Check image size, if over 2500x2500, reject the request
	// and return an error.
	if handler.Size > 4500*4500 {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, `{"status": "error", "message": "Image size is too large"}`)
		return
	}

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
	convertedFile, err := convertImage(buf, format)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, `{"status": "error", "message": "Failed to convert image", "error": "%s"}`, err)
		return
	}

	// Store the converted image in cloud storage
	// and return an authenticated URL to the converted image.
	convertedImageURL, err := storeImage(convertedFile, format)

	// Create and populate the response object
	response := map[string]string{
		"status":      "success",
		"status_code": "200",
		"image_url":   convertedImageURL,
		"message":     "image  successfully converted to " + format + " format",
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

	//

}

// convertImage converts an image to a different format
func convertImage(buf []byte, format string) ([]byte, error) {
	// Use disintegration/imaging to convert the image type
	// and return the converted image as a byte slice.

	// Load the image from byte slice
	src, err := imaging.Decode(bytes.NewReader(buf))
	if err != nil {
		return nil, err
	}

	// Convert the image to the desired format
	// dst := image.NewRGBA(image.Rect(0, 0, src.Bounds().Dx(), src.Bounds().Dy()))
	// draw.Draw(dst, dst.Bounds(), image.NewUniform(color.White), image.ZP, draw.Src)
	// draw.Draw(dst, dst.Bounds(), src, src.Bounds().Min, draw.Over)

	// // Create a buffer to store the converted image
	outputBuf := bytes.NewBuffer(nil)

	// If the image is being converted to WEBP,
	// use chai2010/webp to convert the image to WEBP
	// and return the converted image as a byte slice.
	if format == "webp" {
		// Convert the image to WEBP
		// and return the converted image as a byte slice.
		// Encode lossless webp
		if err = webp.Encode(outputBuf, src, &webp.Options{Lossless: true}); err != nil {
			log.Println(err)
		}
	} else if format == "png" {
		// Save the converted image to the buffer in the desired format
		err = png.Encode(outputBuf, src)
		if err != nil {
			return nil, err
		}
	} else {
		// Save the converted image to the buffer in the desired format
		err = imaging.Encode(outputBuf, src, formatMapping[format])
		if err != nil {
			return nil, err
		}
	}

	return outputBuf.Bytes(), nil
}

// storeImage stores an image in cloud storage and returns an authenticated URL to the image
func storeImage(buf []byte, newExt string) (string, error) {
	// Simulate storing the image in cloud storage
	// Here, we generate a random URL with a timestamp as an example
	// In a real implementation, you would call the appropriate cloud storage service API to store the image

	// Generate a random URL with timestamp
	rand.Seed(time.Now().UnixNano())
	randNum := rand.Intn(1000)
	timestamp := time.Now().Format("20060102150405")
	name := fmt.Sprintf("webchest_%d_image_%s.%s", randNum, timestamp, newExt)

	// Get the context
	ctx := context.Background()

	// Bucket name
	bucketName := "webchest_image_converter"

	// Get the global storage client instance
	client := config.GetStorageClient()

	// Create a new instance of CloudStorageRepository
	cloudStorageRepo, err := repositories.NewCloudStorageRepository(bucketName, client)
	if err != nil {
		return "", err
	}

	// Call UploadImage on the cloudStorageRepo instance to upload the image to cloud storage
	// and return an authenticated URL to the image.
	url, err := cloudStorageRepo.UploadImage(ctx, name, bytes.NewReader(buf))
	if err != nil {
		return "", err
	}

	// Return the authenticated URL to the image
	return url, nil
}

// isFormatSupported checks if a given image format is supported
func isFormatSupported(format string) bool {
	// Supported image formats
	supportedFormats := []string{"webp", "jpg", "jpeg", "png", "bmp", "gif", "tiff"}
	// Check if the format is supported
	for _, f := range supportedFormats {
		if f == format {
			return true
		}
	}
	return false
}
