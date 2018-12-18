FROM golang:alpine AS build

ENV GO111MODULE=on

# RUN mkdir /go/src/app
WORKDIR /go/src/app

LABEL maintainer="github@shanaakh.pro"

RUN apk add bash ca-certificates git gcc g++ libc-dev

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

# RUN go build -o /go/bin/server src/*.go
# RUN CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go install -a -tags netgo  ./src/*.go
RUN CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -o /go/bin/server src/*.go

FROM alpine

COPY --from=build /go/bin/server /app/server

WORKDIR /app

EXPOSE 8080

ENTRYPOINT ["./server"]

