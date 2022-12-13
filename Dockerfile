FROM golang:1.19-alpine as builder
COPY gosrc /go/src/nutcheck
WORKDIR /go/src/nutcheck
RUN apk add --no-cache git gcc libc-dev && \
    go mod tidy && \
    GO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -tags netgo -installsuffix netgo -ldflags '-w' -o nutcheck .


FROM alpine:3.9
LABEL maintainer "benj.saiz@gmail.com"
COPY --from=builder /go/src/nutcheck/nutcheck /usr/bin/nutcheck

CMD ["nutcheck"]
