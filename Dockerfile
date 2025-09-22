FROM golang:1.25.1-alpine

WORKDIR /app/

COPY . /app/

RUN go build -tags netgo -ldflags '-s -w' -o app

EXPOSE 8000

ENTRYPOINT [ "./app" ]
