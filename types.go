package toolkit

import "log"

// JSON response type structure.
type JSONResponse struct {
	Error   bool        `json:"error"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

type Tools struct {
	MAXJSONSize int
	ErrorLog    *log.Logger
	InfoLog     *log.Logger
}
