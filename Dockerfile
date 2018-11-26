FROM golang:1.11

LABEL maintainer="github@shanaakh.pro"

ENV PORT=8080
ENV GOPATH=/go

RUN mkdir /go/src/app
WORKDIR /go/src/app
COPY . /go/src/app

RUN go get ./...
RUN go build -o main src/*.go

EXPOSE $PORT

CMD ["./main"]

