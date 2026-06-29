package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"homelens/server"
	"homelens/server/alert"
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
	_ = godotenv.Load()

	var token, addr string

	token = os.Getenv("HOMELENS_AUTH_TOKEN")
	addr = os.Getenv("HOMELENS_SERVER_ADDR")

	if token == "" || addr == "" {
		log.Fatal("HOMELENS_AUTH_TOKEN and HOMELENS_SERVER_ADDR environment variables must be set")
	}

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	dbb, err := sql.Open("sqlite", "homelens.db?_pragma=journal_mode(WAL)&_pragma=busy_timeout(5000)")
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}

	if _, err := dbb.ExecContext(ctx, db.Schema); err != nil {
		return fmt.Errorf("failed to run migrations: %w", err)
	}

	queries := db.New(dbb)

	agentRegistry := server.NewAgentRegistry()

	alertEngine := alert.NewEngine(queries, agentRegistry)

	go func() {
		err := alertEngine.Start(ctx)
		if err != nil {
			log.Printf("err starting alert engine: %v", err)
			cancel()
		}
	}()

	agentServer := server.NewAgentServer(log.Printf, token, agentRegistry, queries, alertEngine)

	api := api.NewAPI(log.Printf, agentRegistry, queries, alertEngine)

	mux := http.NewServeMux()
	mux.Handle("/ws", agentServer)
	mux.Handle("/", api.ServeFrontend())
	mux.HandleFunc("GET /api/agents", api.GetAgents)
	mux.HandleFunc("GET /api/agents/ws", api.HandleWebsocket)
	mux.HandleFunc("GET /api/agents/{guid}", api.GetSnapshots)
	mux.HandleFunc("POST /api/agents/update-name", api.UpdateAgentName)
	mux.HandleFunc("POST /api/alerts", api.SaveAlertConfig)
	mux.HandleFunc("GET /api/alerts", api.GetAlertConfig)

	serverHTTP := &http.Server{
		Addr:    addr,
		Handler: corsMiddleware(mux),
	}

	go func() {
		log.Printf("Server listening on %s", addr)
		if err := serverHTTP.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Printf("HTTP server error: %v", err)
			cancel()
		}
	}()

	<-ctx.Done()
	log.Println("Shutting down gracefully...")

	return serverHTTP.Close()
}
