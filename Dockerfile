FROM golang:1.16-alpine

EXPOSE 5000

RUN mkdir /app

COPY ["go.mod", "go.sum", "/app/"]

WORKDIR /app

RUN go mod download

COPY * /app/

RUN go build -o callcaptcha-service .

CMD ["./callcaptcha-service"]