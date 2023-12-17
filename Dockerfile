# Use the official Golang image to create a build artifact.
FROM golang:1.21.4 as builder

# Set the Current Working Directory inside the container.
WORKDIR /app

# Copy the Go Mod and Sum files.
COPY go.mod go.sum ./

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed.
RUN go mod download && go mod verify

# Copy the source code from the current directory to the Working Directory inside the container.
COPY . .


# Build the Go app. create a binary executable named `app` in the `/usr/local/bin/` folder.
RUN go build -v -o /usr/local/bin/app cmd/web/*.go

# Expose port 8080 to the outside world.
EXPOSE 8080

# Command to run the executable.
CMD ["app"]

# docker build -t go-booking-app:latest .
# docker run -it -p 8585:8080 --name go-app go-booking-app:latest
