package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"homelens/server"
	"homelens/server/api"
	"homelens/server/db"

	"github.com/joho/godotenv"
	_ "modernc.org/sqlite"
)

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	godotenv.Load()

	var token, addr string

	token = os.Getenv("HOMELENS_AUTH_TOKEN")
	addr = os.Getenv("HOMELENS_SERVER_ADDR")

	if token == "" || addr == "" {
		log.Fatal("HOMELENS_AUTH_TOKEN and HOMELENS_SERVER_ADDR environment variables must be set")
	}

	ctx := context.Background()
	dbb, err := sql.Open("sqlite", "homelens.db")
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}

	if _, err := dbb.ExecContext(ctx, db.Schema); err != nil {
		return fmt.Errorf("failed to run migrations: %w", err)
	}

	queries := db.New(dbb)

	agentRegistry := server.NewAgentRegistry()

	agentServer := server.NewAgentServer(log.Printf, token, agentRegistry, queries)

	api := api.NewAPI(log.Printf, agentRegistry, queries)

	mux := http.NewServeMux()
	mux.Handle("/ws", agentServer)
	mux.HandleFunc("/api/agents", api.GetAgents)
	mux.HandleFunc("/api/snapshots", api.GetSnapshots)

	fmt.Println("Server listening on port 6969")
	return http.ListenAndServe(addr, corsMiddleware(mux))
}
