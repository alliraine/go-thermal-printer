package dto

// PrintImageRequest is the JSON payload for POST /api/v1/printer/print-image
type PrintImageRequest struct {
	ImageBase64  string `json:"imageBase64" form:"imageBase64" binding:"required"`
	MaxWidthDots int    `json:"maxWidthDots,omitempty" form:"maxWidthDots"`
}
