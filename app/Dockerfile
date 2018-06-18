FROM golang:latest

RUN mkdir -p /go/src/app
WORKDIR /go/src/app

ADD . /go/src/app

RUN go get -v

CMD go run main.go

EXPOSE 8080
