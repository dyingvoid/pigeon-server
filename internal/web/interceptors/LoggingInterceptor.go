package interceptors

import (
	"context"
	"google.golang.org/grpc"
	"log"
)

type LoggingInterceptor struct {
	logger *log.Logger
}

func (li *LoggingInterceptor) Intercept(
	ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler) (interface{}, error) {
	log.Printf("Received request for method: %s", info.FullMethod)

	resp, err := handler(ctx, req)

	if err != nil {
		log.Printf("Error handling request for method %s: %v", info.FullMethod, err)
	} else {
		log.Printf("Successfully handled request for method: %s", info.FullMethod)
	}

	return resp, err
}

func NewLoggingInterceptor(logger *log.Logger) *LoggingInterceptor {
	return &LoggingInterceptor{logger}
}
