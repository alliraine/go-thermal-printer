package dto

type PrinterStatusDto struct {
	PrinterStatus         byte `json:"printerStatus"`
	OfflineStatus         byte `json:"offlineStatus"`
	ErrorStatus           byte `json:"errorStatus"`
	ContinuousPaperStatus byte `json:"continuousPaperStatus"`
}

type PrinterPrintDto struct {
	Data string `json:"data" binding:"required"`
}

type PrinterPrintTemplateDto struct {
	TemplateFile string         `json:"templateFile" binding:"required"`
	Variables    map[string]any `json:"variables"`
}
