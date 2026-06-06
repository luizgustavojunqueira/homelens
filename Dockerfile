FROM oven/bun:alpine AS frontend-builder
WORKDIR /app/ui

COPY /ui/package.json ui/bun.lock* ./
RUN bun install

COPY ui/ .
RUN bun run build


FROM golang:1.26-alpine AS backend-builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .

COPY --from=frontend-builder /app/ui/dist ./ui/dist
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o homelens-server ./cmd/server

FROM alpine:3.23
RUN apk --no-cache add ca-certificates tzdata
WORKDIR /app

COPY --from=backend-builder /app/homelens-server .
RUN mkdir -p /app/data

EXPOSE 80
VOLUME /app/data

CMD ["./homelens-server"]
