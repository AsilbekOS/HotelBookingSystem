# Base image
FROM golang:latest

# Konteyner ichida ishchi katalogini o'rnating
WORKDIR /app

# Copy the source code into the container
COPY . .
COPY ./tls/localhost.pem ./tls/localhost.pem
COPY ./tls/localhost-key.pem ./tls/localhost-key.pem

# Build the Go app
RUN go build -o main ./api/cmd/main.go

# Command to run the executable
CMD ["./main"]
