#!/bin/bash

set -e

echo "======================================"
echo "      Running Backend Tests (Go)      "
echo "======================================"
go test ./... -v -race

echo ""
echo "======================================"
echo "  Running Frontend Tests (React/UI)   "
echo "======================================"
cd ui
npm test
cd ..

echo ""
echo "======================================"
echo "    Running Load Tests (k6)           "
echo "======================================"

K6_CMD=""
if command -v k6 &> /dev/null; then
    K6_CMD="k6"
elif [ -f "./k6-v0.49.0-linux-amd64/k6" ]; then
    K6_CMD="./k6-v0.49.0-linux-amd64/k6"
else
    echo "k6 not found. Skipping load tests."
    echo "To run load tests, please install k6 (e.g. 'sudo apt install k6')."
    echo "======================================"
    echo "    Backend and UI tests passed!      "
    echo "======================================"
    exit 0
fi

echo "Starting temporary backend server on port 8080..."

HOMELENS_DB_PATH=/tmp/homelens_test.db \
HOMELENS_AUTH_TOKEN=dev-token \
HOMELENS_SERVER_ADDR=localhost:8080 \
go run ./cmd/server/main.go > /dev/null 2>&1 &
SERVER_PID=$!

sleep 3

echo "Running k6 loadtest script..."
TOKEN=dev-token BASE_URL=ws://localhost:8080/ws $K6_CMD run loadtest.js
K6_EXIT_CODE=$?

echo "Stopping temporary server (PID $SERVER_PID)..."
kill $SERVER_PID
wait $SERVER_PID 2>/dev/null || true
rm -f /tmp/homelens_test.db*

if [ $K6_EXIT_CODE -eq 0 ]; then
    echo ""
    echo "======================================"
    echo "         All tests passed!            "
    echo "======================================"
else
    echo ""
    echo "======================================"
    echo "       k6 load tests failed!          "
    echo "======================================"
    exit $K6_EXIT_CODE
fi
