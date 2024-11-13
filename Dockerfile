# Build stage
FROM golang:1.23-alpine as builder

WORKDIR /app

COPY go.mod .
RUN go mod download

COPY . .

# Build the Go app
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o myapp

# Final stage with a minimal base image
FROM alpine:latest

WORKDIR /root/

# Copy the compiled binary from the builder image
COPY --from=builder /app/myapp .

COPY .env .

RUN ls -al

# Expose the necessary port
EXPOSE 8080

# Command to run the application
CMD ["./myapp"]


