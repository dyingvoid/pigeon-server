package web

import (
	"encoding/json"
	"github.com/dyingvoid/pigeon-server/internal/web/authentication"
	"github.com/dyingvoid/pigeon-server/internal/web/requests"
	"log"
	"net/http"
	"strings"
)

type HttpContext struct {
	auth   *authentication.Authentication
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

	request, body, err := parseJson[requests.NewUserRequest](r)
	if err != nil {
		context.logger.Println(err)
		if strings.HasPrefix(err.Error(), "internal") {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		context.badRequest(w)
		return
	}

	if err = context.auth.ValidateChallenge(request.GetSignature(), body); err != nil {
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}

	// TODO writeJson
	json.NewEncoder(w).Encode("validatedData")
}

func (context *HttpContext) HealthCheckHandler(w http.ResponseWriter, r *http.Request) {
	js := `{"status": "alive"}`
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(js))
}
