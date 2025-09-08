FROM golang:1.22-alpine AS base

FROM base AS build

WORKDIR /app

COPY . .

RUN go build ./cmd/dockerHook

FROM alpine:latest

WORKDIR /app

COPY --from=build /app/dockerHook .

EXPOSE 8080

ENTRYPOINT [ "/app/dockerHook" ]