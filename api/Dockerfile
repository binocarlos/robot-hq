FROM golang:1.24-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -o server .

FROM alpine:3.20
WORKDIR /app
COPY --from=builder /app/server .
EXPOSE 8080
CMD ["./server"]

# --- Dev target ---
FROM golang:1.24-alpine AS dev
RUN go install github.com/air-verse/air@v1.61.7
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
EXPOSE 8080
CMD ["air"]
