FROM golang:1.21.0

WORKDIR /app

COPY . .

RUN go mod tidy
RUN go build -o queue .

CMD ["./queue"]
