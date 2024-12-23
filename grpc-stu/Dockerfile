FROM golang:1.23.4-alpine3.21 AS build

LABEL authors="pipi"

WORKDIR /build

COPY . .

RUN go build -ldflags "-s -w" -o /hello-server ./cmd/server

FROM alpine as release

COPY --from=build /hello-server /hello-server

EXPOSE 50051

ENTRYPOINT ["/hello-server"]