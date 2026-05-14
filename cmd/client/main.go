package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"homelens/client"
	"homelens/shared"

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

	agentClient := client.NewAgentClient(log.Printf, token, agentID)
	agentClient.Connect(addr) // TODO: Handle connection errors and retries

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	out := make(chan shared.SystemInfo)
	go client.Collect(ctx, time.Second, out)

	for {
		select {
		case <-sigChan:
			fmt.Println("Shutting down...")
			cancel()
			agentClient.Disconnect()
			return
		case snapshot := <-out:
			if err := agentClient.SendSnapshot(snapshot); err != nil {
				fmt.Println("send error:", err)
				return
			}
		}
	}
}
