FROM alpine:3.16.2

WORKDIR /app

COPY ./build/static ./static
COPY ./build/server ./server

EXPOSE 8080

CMD ./server
