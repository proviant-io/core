# Dockerfile References: https://docs.docker.com/engine/reference/builder/

# Start from the latest golang base image
FROM golang:latest

# Add Maintainer Info
LABEL maintainer="Grigorii Merkushev <brushknight@gmail.com>"

# Set the Current Working Directory inside the container
WORKDIR /app

# Declare mountpoint for db volume
VOLUME ["/app/db"]

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download

# Copy the source from the current directory to the Working Directory inside the container
COPY cmd ./cmd
COPY internal ./internal

# copy Makefile
COPY Makefile ./Makefile

# Build the Go app
RUN make docker/compile

# Expose port 80 to the outside world
EXPOSE 80

# Command to run the executable
CMD ["./app"]
