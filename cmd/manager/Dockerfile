FROM golang:1.24-alpine as builder

LABEL maintainer="amoCRM"

RUN apk add --no-cache git

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY ../.. .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -a -o manager ./cmd/manager
RUN CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -a -o worker ./cmd/workers


FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

COPY --from=builder /app/manager .
COPY --from=builder /app/worker .
COPY --from=builder /app/.env .

RUN chmod +x /root/manager /root/worker

CMD ["/root/manager", "-c", "2"]
