FROM golang:alpine as builder

WORKDIR /app

COPY . .

RUN go mod download

RUN go build -o videostorage ./cmd/videostorage/main.go

EXPOSE 8080

CMD ["./videostorage"]
