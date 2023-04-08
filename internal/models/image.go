package models

type Image struct {
	ID               string `json:"id"`
	OriginalFileName string `json:"original_file_name"`
	OriginalURL      string `json:"original_url"`
	Width            int    `json:"width"`
	Height           int    `json:"height"`
	Format           string `json:"format"`
	Bytes            int    `json:"bytes"`
	DownloadFilename string `json:"download_filename"`
	DownloadURL      string `json:"download_url"`
}
