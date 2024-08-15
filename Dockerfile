FROM golang:1.16-alpine

WORKDIR /app

COPY . .

RUN go build -o todo-api

FROM alpine:3.12

RUN apk add --no-cache curl tar \
    && curl -OL https://golang.org/dl/go1.22.1.linux-amd64.tar.gz \
    && tar -C /usr/local -xzf go1.22.1.linux-amd64.tar.gz \
    && rm go1.22.1.linux-amd64.tar.gz

ENV PATH="/usr/local/go/bin:$PATH"


EXPOSE 8080

CMD ["/todo-api"]