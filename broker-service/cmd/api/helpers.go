package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
)

type jsonResponse struct {
	OK      bool   `json:"ok"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
}

func (app *Config) readJSON(w http.ResponseWriter, r *http.Request, dst any) error {
	maxBytes := 1_048_576
	r.Body = http.MaxBytesReader(w, r.Body, int64(maxBytes))
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	err := dec.Decode(dst)
	if err != nil {
		var syntaxError *json.SyntaxError
		var unmarshalTypeError *json.UnmarshalTypeError
		switch {
		case errors.As(err, &syntaxError):
			return fmt.Errorf("body contains badly-formed JSON (at position %d)", syntaxError.Offset)
		case errors.Is(err, io.ErrUnexpectedEOF):
			return fmt.Errorf("body contains badly-formed JSON")
		case errors.As(err, &unmarshalTypeError):
			return fmt.Errorf("body contains an invalid value for the %q field (at position %d)", unmarshalTypeError.Field, unmarshalTypeError.Offset)
		case strings.HasPrefix(err.Error(), "json: unknown field "):
			fieldName := strings.TrimPrefix(err.Error(), "json: unknown field ")
			return fmt.Errorf("body contains unknown field %s", fieldName)
		case errors.Is(err, io.EOF):
			return fmt.Errorf("body must not be empty")
		case err.Error() == "http: request body too large":
			return fmt.Errorf("body must not be larger than %d bytes", maxBytes)
		default:
			return fmt.Errorf("body contains badly-formed JSON")
		}
	}
	err = dec.Decode(&struct{}{})
	if err != io.EOF {
		return errors.New("body must only contain a single JSON object")
	}
	return nil
}

func (app *Config) writeJSON(w http.ResponseWriter, status int, data any, headers http.Header) error {
	out, err := json.MarshalIndent(data, "", "\t")
	if err != nil {
		return err
	}
	for key, value := range headers {
		w.Header()[key] = value
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(out)
	return nil
}
