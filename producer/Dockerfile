FROM golang:1.15

RUN mkdir -p /build
WORKDIR /build/
COPY . .
# CGO_ENABLED must be disabled to run go binary in Alpine
RUN CGO_ENABLED=0 GOOS=linux go build -mod vendor -o bin/producer github.com/Ezetowers/nsq-test/producer


FROM busybox:latest
COPY --from=0 /build/bin/producer /producer
ENTRYPOINT ["/bin/sh"]