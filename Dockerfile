FROM golang:1.9.0-alpine3.6
RUN apk add --no-cache git \
    && go get -u github.com/golang/dep/cmd/dep
COPY api /go/src/local/deepsea/api
# COPY vendor /go/src/local/deepsea/vendor
ADD Gopkg.toml /go/src/local/deepsea
ADD Gopkg.lock /go/src/local/deepsea
WORKDIR /go/src/local/deepsea
RUN dep ensure \
    && go build -o deepsea-api local/deepsea/api



FROM alpine:edge
RUN apk add --no-cache upx
COPY --from=0 /go/src/local/deepsea/deepsea-api /deepsea/deepsea-api
RUN upx --brute /deepsea/deepsea-api



FROM alpine:3.6
LABEL maintainer="dino.hensen@gmail.com"
COPY --from=1 /deepsea/deepsea-api /deepsea/deepsea-api
EXPOSE 8080
CMD ["/deepsea/deepsea-api"]