package main

import (
	"first_try/server/sv"
	"flag"
	"log"
)

func main() {
	port := flag.String("port", ":5000", "gRPC server port")
	flag.Parse()

	err := sv.RunServer(*port)
	if err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}
