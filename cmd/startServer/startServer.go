package main

import (
	"context"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/g0shi4ek/grcp_weather/config"
	pb "github.com/g0shi4ek/grcp_weather/gen/go"
	"github.com/g0shi4ek/grcp_weather/internal/server/handlers"
	"github.com/g0shi4ek/grcp_weather/internal/server/repository"
	"github.com/g0shi4ek/grcp_weather/internal/server/services"
	"github.com/g0shi4ek/grcp_weather/pkg/database"

	"google.golang.org/grpc"
)


func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	cfg := config.LoadConfig()

	port := cfg.ServerConf.Port
	
	lis, err := net.Listen("tcp", ":" + port)
	if err != nil {
		return 
	}
	s := grpc.NewServer()

	pool, err := database.NewPool(ctx, cfg)
	if err != nil{
		return
	}
	
	weatherServer := handlers.NewWeatherHandler(
		services.NewweatherService(
			repository.NewWeatherRepo(pool),
		),
		cfg,
	)
	pb.RegisterWeatherServiceServer(s, weatherServer)

	go func() {
		<-ctx.Done()
		log.Println("Shutting down server...")
		s.GracefulStop()
	}()

	log.Printf("gRPC server started on port %s\n", port)
	s.Serve(lis)
}

