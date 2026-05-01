package main

import (
	"fmt"
	"log"
	"net/http"

	"homelens/server"
)

func main() {
	agentServer := server.NewAgentServer(log.Printf)

	http.Handle("/ws", agentServer)
	fmt.Println("Server listening on port 6969")
	http.ListenAndServe(":6969", nil)
}
