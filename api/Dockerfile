FROM golang:1.8

ENV GOPATH /go

FROM golang:1.8
COPY ./api /go/src/api
COPY ./api/main.go /go/src/app/main.go

WORKDIR /go/src/app

RUN go get -d -v ./...
RUN go install -v ./...

CMD [ "app" ]
