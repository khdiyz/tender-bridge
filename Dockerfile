FROM golang:1.22-alpine

WORKDIR /app

# Copy and install dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the source code
COPY . .

# Install PostgreSQL client
RUN apk add --no-cache postgresql-client

# Make wait-for-postgres.sh executable
RUN chmod +x wait-for-postgres.sh

# Build the Go application
RUN go build -o main ./cmd/app/main.go

CMD ["./main"]
