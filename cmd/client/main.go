package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"homelens/client"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("could not load .env file: ", err)
	}

	var token, agentID, addr, secondsInterval string

	token = os.Getenv("HOMELENS_AUTH_TOKEN")
	agentID = os.Getenv("HOMELENS_AGENT_ID")
	addr = os.Getenv("HOMELENS_SERVER_ADDR")
	secondsInterval = os.Getenv("HOMELENS_SECONDS_INTERVAL")

	if token == "" || agentID == "" || addr == "" {
		log.Fatal("HOMELENS_AUTH_TOKEN, HOMELENS_AGENT_ID, and HOMELENS_SERVER_ADDR environment variables must be set")
	}

	interval, err := strconv.Atoi(secondsInterval)
	if err != nil {
		interval = 10
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	ctx, stop := signal.NotifyContext(ctx, os.Interrupt, syscall.SIGTERM)
	defer stop()

	agentClient := client.NewAgentClient(log.Printf, token, agentID, addr)

	if err := agentClient.Run(ctx, time.Second*time.Duration(interval)); err != nil {
		log.Fatalf("agent client error: %v", err)
	}
}
