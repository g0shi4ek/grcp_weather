package main

import (
	"first_try/server"
	"flag"
	"log"
)

func main() {
	port := flag.String("port", ":5000", "gRPC server port")
	flag.Parse()

	err := server.RunServer(*port)
	if err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}
