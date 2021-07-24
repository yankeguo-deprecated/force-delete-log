FROM golang:1.16 AS builder
ENV GOPROXY https://goproxy.io
ENV CGO_ENABLED 0
WORKDIR /go/src/app
ADD . .
RUN go build -mod vendor -o /force-delete-log

FROM alpine:3.12
COPY --from=builder /force-delete-log /force-delete-log
CMD ["/force-delete-log", "/data"]