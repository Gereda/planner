FROM golang:latest

ENV GIN_MODE=release
WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY . .

RUN go build -o /restapi ./cmd

EXPOSE 8000

CMD [ "/restapi"]