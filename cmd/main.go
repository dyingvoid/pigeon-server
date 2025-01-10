package main

import (
	"fmt"
	"github.com/dyingvoid/pigeon-server/internal/mongodb"
	"github.com/dyingvoid/pigeon-server/internal/web/authentication"
	"github.com/dyingvoid/pigeon-server/internal/web/handlers"
	"github.com/dyingvoid/pigeon-server/internal/web/interceptors"
	pb "github.com/dyingvoid/pigeon-server/internal/web/proto"
	"github.com/redis/go-redis/v9"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
	"os"
	"time"
)

type application struct {
	Auth   *authentication.Authentication
	Logger *log.Logger
	Mongo  *mongodb.Database
	Redis  *redis.Client
}

func main() {
	db := NewMongo()
	client := NewRedis()

	logger := log.New(os.Stdout, "INFO: ", log.Ldate|log.Ltime)
	auth := authentication.NewAuthentication(
		client,
		32,
		10*time.Second,
	)

	app := application{
		Logger: logger,
		Auth:   &auth,
		Mongo:  db,
		Redis:  client,
	}

	errChan := make(chan error, 2)
	go NewGRPC(&app, ":50051", errChan)
	for err := range errChan {
		if err != nil {
			log.Fatalf("server error: %+v", err)
		}
	}
}

func NewRedis() *redis.Client {
	// TODO env file
	// TODO what are the options?
	options := &redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
		Protocol: 2,
	}
	return redis.NewClient(options)
}

func NewMongo() *mongodb.Database {
	// TODO env file, or some building steps
	dbConfig := mongodb.MongoConfig{
		ConnectionString:   "mongodb://localhost:27017",
		DatabaseName:       "pigeon",
		UserCollectionName: "users",
	}
	db, err := mongodb.NewMongoDb(dbConfig)
	if err != nil {
		log.Fatal(err)
	}
	return db
}

func NewGRPC(app *application, port string, ch chan<- error) {
	loggingInterceptor := interceptors.NewLoggingInterceptor(app.Logger)

	server := grpc.NewServer(
		grpc.UnaryInterceptor(loggingInterceptor.Intercept),
	)

	userService := handlers.NewUserService(app.Mongo.UserRepository)
	challengeService := handlers.NewChallengeService(app.Auth)

	pb.RegisterUserServiceServer(server, userService)
	pb.RegisterChallengeServiceServer(server, challengeService)

	reflection.Register(server)

	listener, err := net.Listen("tcp", port)
	if err != nil {
		app.Logger.Println("grpc set up error")
		ch <- fmt.Errorf("could not listen on port %s, %w", port, err)
	}

	app.Logger.Println("Starting GRPC server on port " + port)
	if err = server.Serve(listener); err != nil {
		app.Logger.Println("grpc set up error")
		ch <- fmt.Errorf("could not start grpc server on port %s, %w", port, err)
	}

	ch <- nil
}
