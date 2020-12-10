FROM golang:1.13-alpine

WORKDIR /go/app

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .


RUN go build -o server

EXPOSE 8080

CMD ["/go/app/server"]
