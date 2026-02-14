package response

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
)

type SuccessEnvelope struct {
	Data any `json:"data"`
	Meta any `json:"meta,omitempty"`
}

type ErrorEnvelope struct {
	Error APIError `json:"error"`
}

type APIError struct {
	Message string `json:"message"`
	Details any    `json:"details,omitempty"`
}

func WriteJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

func WriteData(w http.ResponseWriter, status int, data any, meta any) {
	WriteJSON(w, status, SuccessEnvelope{
		Data: data,
		Meta: meta,
	})
}

func WriteError(w http.ResponseWriter, status int, message string, details any) {
	WriteJSON(w, status, ErrorEnvelope{
		Error: APIError{
			Message: message,
			Details: details,
		},
	})
}

func DecodeJSON(w http.ResponseWriter, r *http.Request, dst any) error {
	r.Body = http.MaxBytesReader(w, r.Body, 1<<20)
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()

	if err := dec.Decode(dst); err != nil {
		var syntaxError *json.SyntaxError
		var unmarshalTypeError *json.UnmarshalTypeError

		switch {
		case errors.As(err, &syntaxError):
			return errors.New("invalid JSON syntax")
		case errors.Is(err, io.EOF):
			return errors.New("request body is required")
		case errors.As(err, &unmarshalTypeError):
			return errors.New("invalid JSON type for field " + unmarshalTypeError.Field)
		default:
			return err
		}
	}

	// Memastikan tidak ada json kedua
	if err := dec.Decode(&struct{}{}); err != io.EOF {
		return errors.New("request body must contain only one JSON object")
	}

	return nil
}
