package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"homelens/client"

	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()

	var token, agentID, addr string

	token = os.Getenv("HOMELENS_AUTH_TOKEN")
	agentID = os.Getenv("HOMELENS_AGENT_ID")
	addr = os.Getenv("HOMELENS_SERVER_ADDR")

	if token == "" || agentID == "" || addr == "" {
		log.Fatal("HOMELENS_AUTH_TOKEN, HOMELENS_AGENT_ID, and HOMELENS_SERVER_ADDR environment variables must be set")
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	ctx, stop := signal.NotifyContext(ctx, os.Interrupt, syscall.SIGTERM)
	defer stop()

	agentClient := client.NewAgentClient(log.Printf, token, agentID, addr)

	if err := agentClient.Run(ctx, time.Second); err != nil {
		log.Fatalf("agent client error: %v", err)
	}
}
