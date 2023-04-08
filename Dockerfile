FROM golang:1.20-bullseye

RUN mkdir /app

ADD . /app

WORKDIR /app

RUN go mod download
RUN go build -o main ./cmd/main/app.go

EXPOSE 8080

CMD [ "/app/main" ]