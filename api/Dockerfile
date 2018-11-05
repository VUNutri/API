FROM golang:1.8

<<<<<<< HEAD
COPY ./api /go/src/api

WORKDIR /go/src/app
COPY ./api/main.go .
=======
ENV GOPATH /go

FROM golang:1.8
COPY ./api /go/src/api
COPY ./api/main.go /go/src/app/main.go

WORKDIR /go/src/app
>>>>>>> 076ec2de0f29ea5aac36193b29775cf9360e4ec4

RUN go get -d -v ./...
RUN go install -v ./...

<<<<<<< HEAD
CMD [ "app" ]
=======
CMD [ "app" ]
>>>>>>> 076ec2de0f29ea5aac36193b29775cf9360e4ec4
