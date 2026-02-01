FROM golang:alpine3.18 AS builder

#ENV GOPROXY=https://goproxy.cn,direct

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

# Set TimeZone to JST+9
RUN apk add -U tzdata
ENV TZ=Asia/Tokyo
RUN cp /usr/share/zoneinfo/Asia/Tokyo /etc/localtime

WORKDIR /app
COPY bootstrap/src/resources/* /app/conf/

COPY --from=builder /app/myapp .

#RUN apk add --no-cache nginx
#
#COPY UI02_WEB-UI/ito-ems-front/.output/public /usr/share/nginx/html
#
#COPY UI02_WEB-UI/local-ems-ito-go/vue/nginx.conf /etc/nginx/nginx.conf
#
#ENV BEEGO_RUNMODE=prod

EXPOSE 28666

ENTRYPOINT ["/bin/sh", "-c", "/app/myapp"]