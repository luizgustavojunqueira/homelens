package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"homelens/server"

	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()

	var token, addr string

	token = os.Getenv("HOMELENS_AUTH_TOKEN")
	addr = os.Getenv("HOMELENS_SERVER_ADDR")

	if token == "" || addr == "" {
		log.Fatal("HOMELENS_AUTH_TOKEN and HOMELENS_SERVER_ADDR environment variables must be set")
	}

	agentRegistry := server.NewAgentRegistry()

	agentServer := server.NewAgentServer(log.Printf, token, agentRegistry)

	http.Handle("/ws", agentServer)
	fmt.Println("Server listening on port 6969")
	http.ListenAndServe(addr, nil)
}
