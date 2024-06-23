package toolkit

import (
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"os"
)

const defaultMaxSize = 10485760 // 10 MB

// New returns a new toolkit with default values.
func New() Tools {
	return Tools{
		MAXJSONSize: defaultMaxSize,
		InfoLog:     log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime),
		ErrorLog:    log.New(os.Stdout, "Error\t", log.Ldate|log.Ltime),
	}
}

// ReadJSON tries to read the data from request body and converts it from JSON to a variable.
// data arguement needs to be
func (t *Tools) ReadJSON(w http.ResponseWriter, r *http.Request, data interface{}) error {
	maxBytes := 1048576 // 1 megabyte

	r.Body = http.MaxBytesReader(w, r.Body, int64(maxBytes))

	dec := json.NewDecoder(r.Body)
	err := dec.Decode(data)
	if err != nil {
		return err
	}

	err = dec.Decode(&struct{}{})
	if err != io.EOF {
		return errors.New("body must only have a single JSON value")
	}

	return nil
}

func (t *Tools) WriteJSON(w http.ResponseWriter, status int, data interface{}, headers ...http.Header) error {
	out, err := json.Marshal(data)
	if err != nil {
		return err
	}

	if len(headers) > 0 {
		for key, value := range headers[0] {
			w.Header()[key] = value
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	_, err = w.Write(out)
	if err != nil {
		return err
	}

	return nil
}

func (t *Tools) ErrorJSON(w http.ResponseWriter, err error, status ...int) error {
	statusCode := http.StatusBadRequest
	if len(status) > 0 {
		statusCode = status[0]
	}

	payload := JSONResponse{
		Error:   true,
		Message: err.Error(),
	}

	return t.WriteJSON(w, statusCode, payload)
}
