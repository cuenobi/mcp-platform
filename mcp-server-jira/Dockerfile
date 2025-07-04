# Build stage
FROM golang:1.23.9-alpine AS builder

WORKDIR /app

# Copy go mod files
COPY mcp-server-jira/go.mod mcp-server-jira/go.sum ./mcp-server-jira/
COPY shared/proto/gen/go.mod shared/proto/gen/go.sum ./shared/proto/gen/

# Copy shared proto generated files
COPY shared/proto/gen/ ./shared/proto/gen/

# Download dependencies
RUN cd mcp-server-jira && go mod download

# Copy source code
COPY mcp-server-jira/ ./mcp-server-jira/

# Build the application
RUN cd mcp-server-jira && go build -o bin/mcp-server-jira main.go

# Runtime stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copy the binary from builder stage
COPY --from=builder /app/mcp-server-jira/bin/mcp-server-jira .

# Expose port
EXPOSE 50051

CMD ["./mcp-server-jira"] 