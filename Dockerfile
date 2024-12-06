FROM golang:1.21-alpine

RUN apk add --no-cache sqlite sqlite-dev gcc musl-dev

WORKDIR /app

COPY go.* ./
RUN go mod download

COPY . .
RUN go build -o main ./cmd/main.go

ENV PORT=8080
EXPOSE 8080

CMD ["./main"]
