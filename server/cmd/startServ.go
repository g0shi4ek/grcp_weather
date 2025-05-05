package main

import (
	"flag"
	"fmt"
	"log"
	"net"

	pb "github.com/g0shi4ek/grcp_weather/gen/go"
	"github.com/g0shi4ek/grcp_weather/server/internal/handlers"
	"github.com/g0shi4ek/grcp_weather/server/internal/repository"
	"github.com/g0shi4ek/grcp_weather/server/internal/services"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)



func StartServer(port string) error {
	port = "50051"

	lis, err := net.Listen("tcp", ":" + port)
	if err != nil {
		return status.Errorf(codes.Internal, "failed to listen: %v", err)
	}
	s := grpc.NewServer()

	if pool, err := NewPostgresPool(); err != nil{
		return status.Errorf(codes.Internal, "failed to create pool: %v", err)
	}
	
	weatherServer := handlers.NewWeatherHandler(
		services.NewweatherService(
			repository.NewWeatherRepo(pool),
		),
	)
	pb.RegisterWeatherServiceServer(s, weatherServer)
	fmt.Printf("gRPC server started on port %s\n", port)
	return s.Serve(lis)
}

