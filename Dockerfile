FROM golang:1.16

WORKDIR /app

COPY ./ ./

RUN go mod download

RUN go build -o main ./cmd/main.go

EXPOSE 8080

CMD ["./main"]