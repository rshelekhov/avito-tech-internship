# Stage 1: Build the application
FROM golang:1.23-alpine3.19 AS builder

WORKDIR /src

# Setup base software for building an app
RUN apk update && apk add --no-cache ca-certificates git make

# Install golang-migrate
RUN go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

# Download dependencies
COPY go.mod go.sum ./
RUN go mod download -x && go mod verify

# Copy application source
COPY . .

# Build the application
RUN go build -o /app ./cmd/app

# Stage 2: Prepare the final runtime image
FROM golang:1.23-alpine3.19 AS runner

RUN apk update && apk add --no-cache ca-certificates make postgresql-client wget

WORKDIR /src

# Copy go.mod and go.sum for module info
COPY go.mod go.sum ./

# Copy all source files for tests
COPY . .

# Copy binary and other files from builder
COPY --from=builder /app ./
COPY --from=builder /go/bin/migrate /usr/local/bin/migrate

CMD ["./app"]