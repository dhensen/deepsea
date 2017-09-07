FROM golang:1.9.0-alpine3.6
RUN apk add --no-cache git \
    && go get -u github.com/golang/dep/cmd/dep

WORKDIR /go/src/local/deepsea

COPY vendor /go/src/local/deepsea/vendor
COPY api /go/src/local/deepsea/api

# Skip Gopkg and dep ensure for now. Just copy vendor and build.
# TODO: use dep ensure when it has vendor verification
# ADD Gopkg.toml /go/src/local/deepsea
# ADD Gopkg.lock /go/src/local/deepsea
# RUN dep ensure

RUN go build -o deepsea-api local/deepsea/api



FROM alpine:3.6
LABEL maintainer="dino.hensen@gmail.com"
COPY --from=0 /go/src/local/deepsea/deepsea-api /deepsea/deepsea-api
EXPOSE 8080
CMD ["/deepsea/deepsea-api"]