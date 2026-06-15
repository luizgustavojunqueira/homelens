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
	_ = godotenv.Load()

	var token, addr, intervalStr string

	token = os.Getenv("HOMELENS_AUTH_TOKEN")
	addr = os.Getenv("HOMELENS_SERVER_ADDR")
	intervalStr = os.Getenv("HOMELENS_SECONDS_INTERVAL")

	if token == "" || addr == "" {
		log.Fatal("HOMELENS_AUTH_TOKEN and HOMELENS_SERVER_ADDR environment variables must be set")
	}

	interval, err := strconv.Atoi(intervalStr)
	if err != nil || interval <= 0 {
		interval = 10
	}

	machineID, err := client.GetMachineID()
	if err != nil {
		log.Fatalf("could not get agent id: %v", err)
	}

	agentClient := client.NewAgentClient(log.Printf, machineID, token, addr)

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	if err := agentClient.Run(ctx, time.Second*time.Duration(interval)); err != nil {
		log.Fatalf("agent client error: %v", err)
	}
}
