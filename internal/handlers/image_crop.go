package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"github.com/NathanielRand/webchest-image-converter-api/internal/config"
	"github.com/NathanielRand/webchest-image-converter-api/internal/repositories"
	"github.com/chai2010/webp"
	"github.com/disintegration/imaging"
)

type ImageCropResponse struct {
	Status   string `json:"status"`
	Message  string `json:"message"`
	ImageURL string `json:"image_url"`
}

// ImageResizeHandler is a handler for the /image-resize endpoint.
func ImageCropHandler(w http.ResponseWriter, r *http.Request) {
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
		response := ImageCropResponse{
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
		response := ImageCropResponse{
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
		response := ImageCropResponse{
			Status:   "error",
			Message:  "Failed to parse form data: " + err.Error(),
			ImageURL: "",
		}
		// Encode the response object as JSON and write it to the response
		json.NewEncoder(rw).Encode(response)
		return
	}

	// Get the height value from the form data
	height := r.FormValue("height")
	if height == "" {
		// Set the response status code to 400 Bad Request
		rw.WriteHeader(http.StatusBadRequest)
		// Create and populate the response object
		response := ImageCropResponse{
			Status:   "error",
			Message:  "Missing height value",
			ImageURL: "",
		}
		// Encode the response object as JSON and write it to the response
		json.NewEncoder(rw).Encode(response)
		return
	}

	// Get the width value from the form data
	width := r.FormValue("width")
	if width == "" {
		// Set the response status code to 400 Bad Request
		rw.WriteHeader(http.StatusBadRequest)
		// Create and populate the response object
		response := ImageCropResponse{
			Status:   "error",
			Message:  "Missing width value",
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
		response := ImageCropResponse{
			Status:   "error",
			Message:  "Failed to get uploaded file: " + err.Error(),
			ImageURL: "",
		}
		// Encode the response object as JSON and write it to the response
		json.NewEncoder(rw).Encode(response)
		return
	}
	defer file.Close()

	// Get the file format
	format, err := getImageFormat(handler.Filename)

	// Create a buffer to store the file data
	buf := make([]byte, handler.Size)

	// Read the file data into the buffer
	_, err = file.Read(buf)
	if err != nil {
		// Set the response status code to 400 Bad Request
		rw.WriteHeader(http.StatusBadRequest)
		// Create and populate the response object
		response := ImageCropResponse{
			Status:   "error",
			Message:  "Failed to read file into buffer: " + err.Error(),
			ImageURL: "",
		}
		// Encode the response object as JSON and write it to the response
		json.NewEncoder(rw).Encode(response)
		return
	}

	// Resize the image
	resizedFile, format, err := cropImage(buf, format, height, width)
	if err != nil {
		// Set the response status code to 500 Internal Server Error
		rw.WriteHeader(http.StatusInternalServerError)
		// Create and populate the response object
		response := ImageCropResponse{
			Status:   "error",
			Message:  "Failed to resize image: " + err.Error(),
			ImageURL: "",
		}
		// Encode the response object as JSON and write it to the response
		json.NewEncoder(rw).Encode(response)
		return
	}

	// Store the converted image in cloud storage
	// and return an authenticated URL to the converted image.
	convertedImageURL, err := storeCroppedImage(resizedFile, format)
	if err != nil {
		// Set the response status code to 500 Internal Server Error
		rw.WriteHeader(http.StatusInternalServerError)
		// Create and populate the response object
		response := ImageCropResponse{
			Status:   "error",
			Message:  "Failed to store image: " + err.Error(),
			ImageURL: "",
		}
		// Encode the response object as JSON and write it to the response
		json.NewEncoder(rw).Encode(response)
		return
	}

	// Create and populate the response object
	response := ImageCropResponse{
		Status:   "success",
		Message:  "Image successfully resized",
		ImageURL: convertedImageURL,
	}

	// Encode the response object as JSON and write it to the response
	// and return an error if the encoding fails.
	if err := json.NewEncoder(rw).Encode(response); err != nil {
		// Set the response status code to 500 Internal Server Error
		rw.WriteHeader(http.StatusInternalServerError)
		// Create and populate the response object
		response := ImageCropResponse{
			Status:   "error",
			Message:  "Failed to encode response as JSON: " + err.Error(),
			ImageURL: "",
		}
		// Encode the response object as JSON and write it to the response
		json.NewEncoder(rw).Encode(response)
		return
	}
}

// resizeImage converts an image to a different format
func cropImage(buf []byte, format string, height string, width string) ([]byte, string, error) {
	// Use disintegration/imaging to resize the image type
	// and return the resized image as a byte slice.

	// Load the image from byte slice
	src, err := imaging.Decode(bytes.NewReader(buf))
	if err != nil {
		return nil, "", fmt.Errorf("failed to decode image: %w", err)
	}

	// Print the image dimensions
	fmt.Printf("Image dimensions: %dx%d", src.Bounds().Dx(), src.Bounds().Dy())

	// Parse the height as integer
	heightInt, err := strconv.Atoi(height)
	if err != nil {
		return nil, "", fmt.Errorf("failed to parse height: %w", err)
	}

	// Parse the width as integer
	widthInt, err := strconv.Atoi(width)
	if err != nil {
		return nil, "", fmt.Errorf("failed to parse width: %w", err)
	}

	// Resize the image with the provided dimensions
	// resized := imaging.Resize(src, widthInt, heightInt, imaging.Lanczos)

	// Resize and crop the srcImage to fill the 100x100px area.
	// croppedImg := imaging.Fill(src, widthInt, heightInt, imaging.Center, imaging.Lanczos)

	// Crop the original image to 300x300px size using the center anchor.
	croppedImg := imaging.CropAnchor(src, widthInt, heightInt, imaging.Center)

	// Create a buffer to store the converted image
	outputBuf := bytes.NewBuffer(nil)

	// Check the image format and encode
	// the image to the buffer with the appropriate format.
	if format == "webp" {
		// Encode the image using chai2010/webp
		if err = webp.Encode(outputBuf, croppedImg, &webp.Options{Lossless: true}); err != nil {
			return nil, "", fmt.Errorf("failed to encode image to WEBP: %w", err)
		}
	} else {
		// Encode the image using disintegration/imaging
		err = imaging.Encode(outputBuf, croppedImg, formatMapping[format])
		if err != nil {
			return nil, "", fmt.Errorf("failed to encode image to %s: %w", format, err)
		}
	}

	// Return the resized image as a byte slice
	return outputBuf.Bytes(), format, nil
}

// storeImage stores an image in cloud storage and returns an authenticated URL to the image
func storeCroppedImage(buf []byte, newExt string) (string, error) {
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
	bucketName := "webchest_image_api_bucket"

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
