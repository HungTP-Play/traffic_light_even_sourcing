FROM golang:1.19-bullseye

WORKDIR /app

# COPY all files in ./event_store to /app
COPY ./event_store/. .

COPY go.mod .

RUN go mod download

RUN ls