FROM golang:alpine3.18 AS builder

ENV GOPROXY=https://goproxy.cn,direct

RUN apk add --no-cache git gcc g++ musl-dev

WORKDIR /app

COPY go.mod go.mod
COPY go.sum go.sum
RUN go mod download

COPY ./ .

RUN CGO_ENABLED=1
RUN go env -w GOCACHE=/go-cache
RUN --mount=type=cache,target=/go-cache GOOS=linux CC=gcc go build -o myapp .

FROM alpine:3.18

WORKDIR /app
COPY src/main/resources/* /app/conf/

COPY --from=builder /app/myapp .

EXPOSE 26666

ENTRYPOINT ["/bin/sh", "-c", "/app/myapp"]