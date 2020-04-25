FROM golang:1.14.1-stretch as builder

COPY ./ /go

RUN cd cmd/simple-http-blackbox-exporter \
  && go get \
  && go build

FROM gcr.io/distroless/cc

COPY --from=builder /go/cmd/simple-http-blackbox-exporter/simple-http-blackbox-exporter /usr/sbin/

ENTRYPOINT ["simple-http-blackbox-exporter"]
