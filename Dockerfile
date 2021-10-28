FROM golang:1.16-alpine as builder
# Set the working directory to /app
WORKDIR /nebula-http-gateway
# Copy the current directory contents into the container at /app
COPY . /nebula-http-gateway
 # Make port available to the world outside this container
ENV GOPROXY https://goproxy.cn

RUN apk --update add \
    go \
    musl-dev \
    util-linux-dev

RUN CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build

FROM alpine

WORKDIR /root
COPY --from=builder ./nebula-http-gateway .
COPY ./conf ./conf

EXPOSE 8080

ENTRYPOINT ["./nebula-http-gateway"]
