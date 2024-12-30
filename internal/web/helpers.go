package web

import (
	"encoding/json"
	"net/http"
)

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
