package handlers

import (
	"fmt"

	"github.com/disintegration/imaging"
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

// getImageFormat returns the image format of the file
func getImageFormat(filename string) (string, error) {
	// Get the file extension
	extension := filename[len(filename)-3:]

	// Check the file extension
	switch extension {
	case "jpg":
		return "jpg", nil
	case "peg":
		return "jpg", nil
	case "png":
		return "png", nil
	case "ebp":
		return "webp", nil
	case "bmp":
		return "bmp", nil
	case "gif":
		return "gif", nil
	case "iff":
		return "tiff", nil
	default:
		return "", fmt.Errorf("Invalid file extension: %s", extension)
	}
}
