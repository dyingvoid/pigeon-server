package web

import (
	"encoding/json"
	"github.com/dyingvoid/pigeon-server/internal/application/models"
	"log"
	"net/http"
)

type HttpContext struct {
	logger *log.Logger
}

type envelope map[string]any

func NewHttpContext(logger *log.Logger) *HttpContext {
	return &HttpContext{logger: logger}
}

func (context *HttpContext) AddUserHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		context.logger.Println("Method not allowed")
		context.methodNotAllowed(w)
	}

	var signedUser models.SignedUser
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&signedUser); err != nil {
		context.logger.Println(err)
		context.badRequest(w)
	}

	// TODO validation logic...

	// TODO writeJson
	json.NewEncoder(w).Encode("validatedData")
}

func (context *HttpContext) HealthCheckHandler(w http.ResponseWriter, r *http.Request) {
	js := `{"status": "alive"}`
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(js))
}
