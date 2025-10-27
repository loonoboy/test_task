FROM golang:1.24-alpine as builder

LABEL maintainer="amoCRM"

RUN apk add --no-cache git

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -a -o main ./cmd/app/main.go

FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

COPY --from=builder /app/main .

COPY --from=builder /app/.env .

EXPOSE 8080

CMD ["./main"]