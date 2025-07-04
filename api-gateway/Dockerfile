FROM golang:1.23.9-alpine AS builder

WORKDIR /app

COPY api-gateway/go.mod api-gateway/go.sum ./api-gateway/

RUN cd api-gateway && go mod download

COPY api-gateway/ ./api-gateway/

RUN cd api-gateway && go build -o bin/api-gateway main.go

FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

COPY --from=builder /app/api-gateway/bin/api-gateway .

EXPOSE 8080

CMD ["./api-gateway"] 