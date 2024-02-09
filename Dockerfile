FROM golang:1.22-alpine

WORKDIR /app

COPY . .

RUN go build ./cmd/dockerHook

EXPOSE 8080

ENTRYPOINT [ "/dockerHook" ]