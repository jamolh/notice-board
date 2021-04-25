FROM golang:1.16.1-alpine as builder

RUN apk add --no-cache libc6-compat
RUN apk add --no-cache git

ENV GO111MODULE=on

WORKDIR /app

COPY .env .
COPY go.mod .
COPY go.sum .

RUN go mod download
RUN go mod verify                                                                                                                                                                                                                               

COPY . . 

RUN go get -u github.com/swaggo/swag/cmd/swag

RUN swag init

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o notice-board .

EXPOSE 50001

CMD ["sh"]

ENTRYPOINT ["/app/notice-board"]