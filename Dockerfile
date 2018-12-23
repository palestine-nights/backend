FROM golang:alpine AS build

ENV GO111MODULE=on

WORKDIR /go/src/app

LABEL maintainer="github@shanaakh.pro"

RUN apk add bash ca-certificates git gcc g++ libc-dev

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

RUN CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -o /go/bin/server src/main.go

FROM alpine

COPY --from=build /go/bin/server /app/server
COPY --from=build /go/src/app/html /app/html
COPY --from=build /go/src/app/templates /app/templates

WORKDIR /app

EXPOSE 8080

ENTRYPOINT ["./server"]

