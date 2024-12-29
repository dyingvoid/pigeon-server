package main

import (
	"encoding/json"
	"fmt"
	"github.com/dyingvoid/pigeon-server/internal"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"os"
)

const version = "dev#0"

type config struct {
	port string
}

type application struct {
	config config
	logger *log.Logger
	db     *internal.Db
}

const get = "GET"

func main() {
	logger := log.New(os.Stdout, "INFO: ", log.Ldate|log.Ltime)

	db := internal.Db{
		Users: make(map[string]internal.User),
	}

	for i := 0; i < 10; i++ {
		db.AddUser(internal.NewUser(10))
	}

	app := application{
		config: config{
			port: ":8080",
		},
		logger: logger,
		db:     &db,
	}

	router := mux.NewRouter()
	router.HandleFunc("/healthcheck", app.healthcheckHandler).Methods(get)
	router.HandleFunc("/users", app.getAllUsersHandler).Methods(get)

	log.Println("Starting server on port " + app.config.port)
	err := http.ListenAndServe(app.config.port, router)
	if err != nil {
		log.Fatal(err)
	}
}

func (app *application) healthcheckHandler(w http.ResponseWriter, r *http.Request) {
	js := `{"status": "available", "version": %q}`
	js = fmt.Sprintf(js, version)

	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(js))
}

func (app *application) getAllUsersHandler(w http.ResponseWriter, r *http.Request) {
	users := app.db.GetAllUsers()

	w.Header().Set("Content-Type", "application/json")

	js, err := json.Marshal(users)
	if err != nil {
		app.logger.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	js = append(js, '\n')
	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}
