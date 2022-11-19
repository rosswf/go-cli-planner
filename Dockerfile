FROM golang:1.19

WORKDIR /usr/src/app

COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .
RUN go build -v -o /usr/local/bin/app cmd/web_server/main.go

CMD ["app"]
EXPOSE 5000