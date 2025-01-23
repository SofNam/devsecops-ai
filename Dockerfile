# Use a specific golang version
FROM golang:latest

WORKDIR /app

COPY . .

RUN go mod tidy
RUN go build -o main ./cmd/scanner/main.go

CMD ["./main"]