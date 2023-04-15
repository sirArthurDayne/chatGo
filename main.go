package main

import (
	"flag"

	"github.com/sirArthurDayne/chatGo/server"
)

const (
	DEFAULT_PORT = 8080
	DEFAULT_HOST = "localhost"
)

var (
	host = flag.String("h", DEFAULT_HOST, "Host to be connected.(default=localhost)")
	port = flag.Int("p", DEFAULT_PORT, "Port to connect.(default=8080)")
)

func main() {
	flag.Parse()
	server := server.NewServer()
	server.LoadServerComponents()
	server.Start(host, port)
}
