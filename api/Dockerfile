FROM golang:1.8

COPY ./api /go/src/api

WORKDIR /go/src/app
COPY ./api/main.go .

RUN go get github.com/pilu/fresh
RUN go-wrapper download
RUN go-wrapper install

CMD [ "app" ]
