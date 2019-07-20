FROM golang:1.8

WORKDIR /go/src/app
COPY . .
RUN go get "github.com/emersion/go-smtp"
RUN go get github.com/streadway/amqp
RUN ["go", "build"]

CMD ["/go/src/app/app"]