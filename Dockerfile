FROM golang:alpine AS builder
WORKDIR /app
COPY . .
RUN apk add build-base && go build -o forum cmd/main.go
FROM alpine:latest
LABEL Authors="Myrzakhmet&Manara"
WORKDIR /app
COPY --from=builder /app .
CMD ["./forum"]