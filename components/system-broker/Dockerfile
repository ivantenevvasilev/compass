FROM --platform=$BUILDPLATFORM golang:1.21.1-alpine3.18 as builder

ENV BASE_APP_DIR /go/src/github.com/kyma-incubator/compass/components/system-broker
WORKDIR ${BASE_APP_DIR}

COPY go.mod go.sum ${BASE_APP_DIR}/
RUN go mod download -x

COPY . ${BASE_APP_DIR}

ARG TARGETOS TARGETARCH
RUN GOOS=$TARGETOS GOARCH=$TARGETARCH go build -v -o /app/system-broker ./cmd/main.go

FROM alpine:3.18.4
LABEL source = git@github.com:kyma-incubator/compass.git
WORKDIR /app

RUN apk --no-cache add curl ca-certificates

COPY --from=builder /app /app

CMD ["/app/system-broker"]
