FROM golang:1.25-alpine AS builder

WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o miniflux-mcp .

FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

COPY --from=builder /app/miniflux-mcp .

# Default port for Streamable HTTP mode. Documentary only - actual listen
# address is controlled by the MCP_HTTP_PORT environment variable at runtime.
# When MCP_HTTP_PORT is unset, the server falls back to stdio mode.
EXPOSE 3000

CMD ["./miniflux-mcp"]
