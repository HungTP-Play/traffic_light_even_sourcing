FROM golang:1.19-bullseye

WORKDIR /app

COPY go.mod ./

RUN go mod download

COPY . .

RUN go build -o /app/bin/event_store .

CMD ["/app/bin/event_store"]