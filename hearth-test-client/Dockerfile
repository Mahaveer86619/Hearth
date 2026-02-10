FROM golang:1.21-alpine

WORKDIR /app

COPY go.mod ./
# COPY go.sum ./ # No dependencies yet
# RUN go mod download

COPY . .

RUN go build -o test-client main.go

EXPOSE 3000

CMD ["./test-client"]
