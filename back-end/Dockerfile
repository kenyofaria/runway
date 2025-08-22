# Stage 1: Build the Go application
FROM golang:1.25.0 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

# Copy the .env file and the rest of the source code
COPY .env .
COPY . .

# Build the application binary
RUN CGO_ENABLED=0 GOOS=linux go build -o main .

# Stage 2: Create a minimal final image
FROM alpine:latest

# Create a non-root user
RUN adduser -D -s /bin/sh appuser

# Set the working directory, which will contain the binary and .env
WORKDIR /app

# Copy the built binary and change its ownership to appuser
COPY --from=builder --chown=appuser:appuser /app/main .
COPY --from=builder --chown=appuser:appuser /app/.env .

RUN mkdir /app/data
RUN touch /app/data/apps.json
RUN touch /app/data/reviews.json
RUN mkdir /app/logs

RUN chown -R appuser:appuser /app/data
RUN chown -R appuser:appuser /app/logs

# Switch to the non-root user
USER appuser

# Expose the application port
EXPOSE 8080

# Command to run the executable
CMD ["./main"]