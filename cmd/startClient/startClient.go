package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/g0shi4ek/grcp_weather/config"
	"github.com/g0shi4ek/grcp_weather/internal/client/app"
	"github.com/g0shi4ek/grcp_weather/internal/client/cli"
)


func main(){
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()
	
	cfg := config.LoadConfig()

	grpcClient, err := app.NewGrpcWeatherClient(cfg)
	if err != nil{
		log.Fatal("Cannot start client")
		return 
	}
	weatherCli := cli.NewWeatherCli(grpcClient, cfg)

	weatherCli.RunApp(ctx)
}


