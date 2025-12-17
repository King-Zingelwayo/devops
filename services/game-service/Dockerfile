# Multi-stage build for minimal image size
FROM golang:1.21-alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main ./cmd/main.go

# Final stage - distroless for security and size
FROM gcr.io/distroless/static-debian11

WORKDIR /

COPY --from=builder /app/main .

EXPOSE 8080

USER nonroot:nonroot

ENTRYPOINT ["./main"]