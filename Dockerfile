# Build stage: Compile the Go application
FROM golang:1.23 AS build-stage

# Set the working directory inside the container
WORKDIR /chat_gpg

# Copy the Go module files to download dependencies
COPY go.mod go.sum ./ 
RUN go mod download

# Copy the rest of the application source code
COPY . ./

# Build the Go application binary
RUN CGO_ENABLED=1 GOOS=linux go build -o /chat_gpg/cmd/chat_gpg ./cmd

# Release stage: Use a minimal base image for deployment
FROM debian:bookworm AS build-release-stage

# Set the working directory inside the container
WORKDIR /

# Create a non-root user
RUN groupadd -g 1000 nonroot && useradd -m -u 1000 -g nonroot nonroot

# Copy the application binary from the build stage
COPY --from=build-stage /chat_gpg/cmd/chat_gpg /chat_gpg

# Create directories for the UI and DB before copying
RUN mkdir -p /chat_gpg/ui /chat_gpg/db

# Copy the `ui` and `db` directories into the image
COPY --from=build-stage /chat_gpg/ui/ /chat_gpg/ui
COPY --from=build-stage /chat_gpg/db /chat_gpg/db
COPY .env /.env

# Expose the port the application will run on
EXPOSE 8080

# Run the application as a non-root user
USER nonroot:nonroot

# Define the entry point for the container
ENTRYPOINT ["/chat_gpg/chat_gpg"]
