FROM golang:1.8

COPY ./api /go/src/api

WORKDIR /go/src/app
COPY ./api/main.go .

RUN go get -d -v ./...
RUN go install -v ./...

CMD [ "app" ]
