# syntax=docker/dockerfile:1

FROM alpine:3.16.2

WORKDIR /app

COPY ./build/static ./static
COPY ./build/server ./server

ENV PRODUCTION=true

EXPOSE 8080

CMD ./server
