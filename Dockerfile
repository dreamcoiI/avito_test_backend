# Используем образ Golang для сборки приложения
FROM golang:1.21-alpine AS builder

WORKDIR /app
COPY . .

RUN go mod download
RUN go build -o main cmd/segmentation/main.go

FROM alpine:latest

RUN apk add bash
COPY --from=builder /app/main /main
EXPOSE 8080

CMD ["./main"]