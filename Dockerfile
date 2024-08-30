FROM golang:1.22.4-alpine3.20 as builder

WORKDIR /app

RUN apk add --no-cache make
RUN go install github.com/swaggo/swag/cmd/swag@latest

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN make generate-swagger

RUN CGO_ENABLED=0 GOOS=linux go build -v -o server ./cmd/server


FROM alpine:latest as runner

WORKDIR /app/

COPY --from=builder /app/server .

EXPOSE 8080

CMD ["./server"]