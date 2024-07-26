FROM golang:latest

WORKDIR /app

ENV GO111MODULE=on

RUN apt install make

COPY go.mod go.sum ./
RUN go mod download

CMD [ "go", "run", "cmd/main.go" ]