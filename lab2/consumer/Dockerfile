FROM golang:1.20.2-alpine3.17

WORKDIR /app

COPY go.mod ./
COPY go.sum ./

RUN go mod download

COPY *.go ./

RUN apk update
RUN apk add gcc musl-dev

ENV CGO_ENABLED=1

RUN go build -o /go-kafka-consumer -tags musl

EXPOSE 5656

CMD [ "/go-kafka-consumer" ]
