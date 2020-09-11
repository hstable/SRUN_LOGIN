FROM golang:alpine AS builder
ADD . /build/
WORKDIR /build/
ENV GO111MODULE=on
ENV GOPROXY=https://goproxy.io
RUN go build -ldflags="-s -w" -o srun_login .

FROM alpine:latest
COPY --from=builder /build/srun_login /usr/bin
ENTRYPOINT ["srun_login"]
