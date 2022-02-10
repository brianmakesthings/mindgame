# syntax=docker/dockerfile:1

FROM golang:1.16-alpine

WORKDIR /app

# ENV GOPATH=/app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

# RUN go get github.com/streadway/amqp github.com/google/uuid

COPY mindgame ./
COPY *.go ./

# RUN go build -o /docker-gs-ping

# EXPOSE 8080

CMD [ "go", "run", "./amqp_server.go" ]