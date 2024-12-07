# Build Stage
FROM golang:1.23-alpine AS build

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o oauth-app ./cmd

# Final Stage
FROM alpine:latest

WORKDIR /root/

# Copy only the compiled binary
COPY --from=build /app/oauth-app .

# Expose the application port
EXPOSE 3000

CMD ["./oauth-app"]