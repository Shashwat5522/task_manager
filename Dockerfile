FROM golang:1.24-alpine

RUN apk add --no-cache git gcc musl-dev ca-certificates

WORKDIR /app

COPY . .

RUN go mod download && go build -o api ./cmd/api

EXPOSE 8080

CMD ["./api"] 
