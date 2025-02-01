FROM golang:1.23.2 AS builder

WORKDIR /server

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o bin/main cmd/main.go

FROM alpine:3.18

WORKDIR /server

COPY --from=builder /server .

RUN chmod +x /server/bin/main

CMD ["/server/bin/main"]