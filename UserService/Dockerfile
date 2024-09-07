# # Base image
# FROM golang:latest

# # Konteyner ichida ishchi katalogini o'rnating
# WORKDIR /app

# # Copy the source code into the container
# COPY . .

# # Build the Go app
# RUN go build -o main ./cmd/main.go

# # Command to run the executable
# CMD ["./main"]


FROM golang:1.22.4 AS builder
WORKDIR /app
COPY go.mod ./
COPY go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o myapp ./cmd

FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/myapp .

COPY .env .
EXPOSE 50051
CMD ["./myapp"]