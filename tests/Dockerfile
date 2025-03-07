FROM --platform=$BUILDPLATFORM golang:1.21.1-alpine3.18 as builder

ENV BASE_TEST_DIR /go/src/github.com/kyma-incubator/compass/tests
WORKDIR ${BASE_TEST_DIR}

COPY . ${BASE_TEST_DIR}

ARG TARGETOS TARGETARCH
RUN GOOS=$TARGETOS GOARCH=$TARGETARCH CGO_ENABLED=0 go test -c ./connectivity-adapter/tests -o connectivity-adapter.test && \
    GOOS=$TARGETOS GOARCH=$TARGETARCH CGO_ENABLED=0 go test -c ./connector/tests -o connector.test && \
    GOOS=$TARGETOS GOARCH=$TARGETARCH CGO_ENABLED=0 go test -c ./director/tests -o director.test && \
    GOOS=$TARGETOS GOARCH=$TARGETARCH CGO_ENABLED=0 go test -c ./director/tests/application -o director-application.test && \
    GOOS=$TARGETOS GOARCH=$TARGETARCH CGO_ENABLED=0 go test -c ./director/tests/notifications -o director-notifications.test && \
    GOOS=$TARGETOS GOARCH=$TARGETARCH CGO_ENABLED=0 go test -c ./director/tests/runtime -o director-runtime.test && \
    GOOS=$TARGETOS GOARCH=$TARGETARCH CGO_ENABLED=0 go test -c ./director/tests/formation -o director-formation.test && \
    GOOS=$TARGETOS GOARCH=$TARGETARCH CGO_ENABLED=0 go test -c ./ns-adapter/tests -o ns-adapter.test && \
    GOOS=$TARGETOS GOARCH=$TARGETARCH CGO_ENABLED=0 go test -c ./external-services-mock/tests -o external-services-mock.test && \
    GOOS=$TARGETOS GOARCH=$TARGETARCH CGO_ENABLED=0 go test -c ./gateway/tests -o gateway.test && \
    GOOS=$TARGETOS GOARCH=$TARGETARCH CGO_ENABLED=0 go test -c ./ord-aggregator/tests -o ord-aggregator.test && \
    GOOS=$TARGETOS GOARCH=$TARGETARCH CGO_ENABLED=0 go test -c ./ord-service/tests -o ord-service.test && \
    GOOS=$TARGETOS GOARCH=$TARGETARCH CGO_ENABLED=0 go test -c ./system-broker/tests -o system-broker.test && \
    GOOS=$TARGETOS GOARCH=$TARGETARCH CGO_ENABLED=0 go test -c ./tenant-fetcher/tests -o tenant-fetcher.test && \
    GOOS=$TARGETOS GOARCH=$TARGETARCH CGO_ENABLED=0 go test -c ./system-fetcher/tests -o system-fetcher.test && \
    GOOS=$TARGETOS GOARCH=$TARGETARCH CGO_ENABLED=0 go test -c ./istio/tests -o istio.test && \
    GOOS=$TARGETOS GOARCH=$TARGETARCH CGO_ENABLED=0 go test -c ./director/bench -o director.bench && \
    GOOS=$TARGETOS GOARCH=$TARGETARCH CGO_ENABLED=0 go test -c ./ord-service/bench -o ord-service.bench && \
    GOOS=$TARGETOS GOARCH=$TARGETARCH CGO_ENABLED=0 go test -c ./pairing-adapter/tests -o pairing-adapter.test

FROM alpine:3.18.4

RUN apk add --no-cache curl

LABEL source=git@github.com:kyma-incubator/compass.git

COPY --from=builder /go/src/github.com/kyma-incubator/compass/tests/connectivity-adapter.test .
COPY --from=builder /go/src/github.com/kyma-incubator/compass/tests/connector.test .
COPY --from=builder /go/src/github.com/kyma-incubator/compass/tests/director.test .
COPY --from=builder /go/src/github.com/kyma-incubator/compass/tests/director-application.test .
COPY --from=builder /go/src/github.com/kyma-incubator/compass/tests/director-runtime.test .
COPY --from=builder /go/src/github.com/kyma-incubator/compass/tests/director-formation.test .
COPY --from=builder /go/src/github.com/kyma-incubator/compass/tests/director-notifications.test .
COPY --from=builder /go/src/github.com/kyma-incubator/compass/tests/ns-adapter.test .
COPY --from=builder /go/src/github.com/kyma-incubator/compass/tests/external-services-mock.test .
COPY --from=builder /go/src/github.com/kyma-incubator/compass/tests/gateway.test .
COPY --from=builder /go/src/github.com/kyma-incubator/compass/tests/ord-aggregator.test .
COPY --from=builder /go/src/github.com/kyma-incubator/compass/tests/ord-service.test .
COPY --from=builder /go/src/github.com/kyma-incubator/compass/tests/system-broker.test .
COPY --from=builder /go/src/github.com/kyma-incubator/compass/tests/tenant-fetcher.test .
COPY --from=builder /go/src/github.com/kyma-incubator/compass/tests/system-fetcher.test .
COPY --from=builder /go/src/github.com/kyma-incubator/compass/tests/istio.test .
COPY --from=builder /go/src/github.com/kyma-incubator/compass/tests/pairing-adapter.test .

COPY --from=builder /go/src/github.com/kyma-incubator/compass/tests/director.bench .
COPY --from=builder /go/src/github.com/kyma-incubator/compass/tests/ord-service.bench .
