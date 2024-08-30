FROM golang:1.22.4-alpine3.20 as builder

WORKDIR /tests

COPY go.mod go.sum ./
RUN go mod download

COPY . .

CMD CGO_ENABLED=0 GOOS=linux go test -v ./...
