FROM golang:1.23.9-alpine AS builder

WORKDIR /app

COPY mcphost/go.mod mcphost/go.sum ./mcphost/
COPY shared/proto/gen/go.mod shared/proto/gen/go.sum ./shared/proto/gen/

COPY shared/proto/gen/ ./shared/proto/gen/

RUN cd mcphost && go mod download

COPY mcphost/ ./mcphost/

RUN cd mcphost && go build -o bin/mcphost main.go

FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

COPY --from=builder /app/mcphost/bin/mcphost .

CMD ["tail", "-f", "/dev/null"] 