FROM golang:alpine3.17 AS build
WORKDIR /app
COPY . .
RUN apk --update-cache upgrade && \
    ls && \
    go get ./internal && \
    go build -o /ctfd internal/*

FROM alpine:3.17.2
WORKDIR /app
RUN apk --update-cache upgrade
COPY --from=build /ctfd .
ENTRYPOINT [ "/app/ctfd" ]