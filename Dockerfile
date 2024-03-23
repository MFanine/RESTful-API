FROM golang:1.18-alpine as builder


# Define build env
ENV GOOS linux
ENV CGO_ENABLED 0

# Add a work directory
WORKDIR /build

# Cache and install dependencies
COPY . .
RUN go mod download

# Build the application
RUN go build -o main .

# Move to /dist directory as the place for resulting binaries
WORKDIR /dist

# Copy binary from build to main folder
RUN cp /build/main .

# Build a small image
FROM scratch

COPY --from=builder /dist/main /
COPY .env /

EXPOSE 8080

# Command to run
ENTRYPOINT ["/main"]
