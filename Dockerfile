# Step 1: Build the Go application in an Alpine environment
FROM golang:1.21.0-alpine3.18 as builder

# Set the working directory in the builder
WORKDIR /app

# Copy the Go mod and sum files to fetch the dependencies
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the entire content to the working directory
COPY . .

# Build the Go app with static linking and strip the binary
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o presence .

# Step 2: Use a 'scratch' image, which is an empty image
FROM scratch

# Copy the compiled binary from the builder to the current container
COPY --from=builder /app/presence /presence

# Set the binary as the default command
ENTRYPOINT ["/presence"]
