FROM golang:1.19-bullseye

WORKDIR /app

COPY ./projector/. .

COPY go.mod .

RUN go mod download