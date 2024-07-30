FROM golang:1.21.1-alpine

WORKDIR /usr/src

COPY go.mod go.sum ./
RUN go mod download

COPY . .

CMD ["./entrypoint.sh"]