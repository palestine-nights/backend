FROM golang:alpine AS build

LABEL maintainer="github@shanaakh.pro"

ENV GOPATH=/go

RUN mkdir /go/src/app
WORKDIR /go/src/app
COPY . /go/src/app

RUN apk update && apk add --no-cache git
RUN go get -v ./...
RUN go build -o /go/bin/server src/*.go

FROM alpine

COPY --from=build /go/bin/server /app/server

WORKDIR /app

EXPOSE 8080

ENTRYPOINT ["./server"]

