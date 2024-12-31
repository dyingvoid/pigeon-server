package web

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// TODO custom errors perhaps
func parseJson[T any](r *http.Request) (*T, []byte, error) {
	defer r.Body.Close()
	body, err := io.ReadAll(r.Body)
	if err != nil {
		return nil, body, fmt.Errorf("internal: %w", err)
	}

	var request T
	decoder := json.NewDecoder(bytes.NewReader(body))
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(&request); err != nil {
		return nil, body, fmt.Errorf("bad request: %w", err)
	}

	return &request, body, nil
}

func (context *HttpContext) writeJson(
	w http.ResponseWriter, status int, data any, headers http.Header) error {
	js, err := json.Marshal(data)
	if err != nil {
		return err
	}

	js = append(js, '\n')
	for k, v := range headers {
		w.Header()[k] = v
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	// TODO, learn why error is ignored here
	w.Write(js)
	return nil
}

func (context *HttpContext) okResponse(w http.ResponseWriter, message any) {
	err := context.writeJson(w, http.StatusOK, message, nil)
	if err != nil {
		context.logger.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func (context *HttpContext) errorResponse(
	w http.ResponseWriter, status int, message any) {
	env := envelope{"error": message}
	err := context.writeJson(w, status, env, nil)
	if err != nil {
		context.logger.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func (context *HttpContext) methodNotAllowed(w http.ResponseWriter) {
	http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
}

func (context *HttpContext) badRequest(w http.ResponseWriter) {
	http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
}

func (context *HttpContext) unauthorized(w http.ResponseWriter) {
	http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
}
