package utils

import (
	"encoding/json"
	"net/http"
)

// SuccessResponse merepresentasikan response standar untuk operasi sukses
type SuccessResponse struct {
	Success bool `json:"success"`
	Data    any  `json:"data"`
	Meta    any  `json:"meta,omitempty"`
}

// ErrorDetail merepresentasikan detail error pada response error
type ErrorDetail struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Details any    `json:"details,omitempty"`
}

// ErrorResponse merepresentasikan response standar untuk operasi gagal
type ErrorResponse struct {
	Success bool        `json:"success"`
	Error   ErrorDetail `json:"error"`
}

// JSONSuccess menulis response sukses berformat JSON sesuai standar skill api-design
func JSONSuccess(w http.ResponseWriter, statusCode int, data any, meta ...any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	var metaData any
	if len(meta) > 0 {
		metaData = meta[0]
	}

	resp := SuccessResponse{
		Success: true,
		Data:    data,
		Meta:    metaData,
	}

	_ = json.NewEncoder(w).Encode(resp)
}

// JSONError menulis response error berformat JSON sesuai standar skill api-design
func JSONError(w http.ResponseWriter, statusCode int, code, message string, details ...any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	var detailData any
	if len(details) > 0 {
		detailData = details[0]
	}

	resp := ErrorResponse{
		Success: false,
		Error: ErrorDetail{
			Code:    code,
			Message: message,
			Details: detailData,
		},
	}

	_ = json.NewEncoder(w).Encode(resp)
}
