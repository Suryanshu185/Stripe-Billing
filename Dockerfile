# Dockerfile
FROM golang:1.21-alpine

WORKDIR /app

RUN apk update && apk add git

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY . .

RUN go build -o server .

EXPOSE 8080
CMD ["./server"]
