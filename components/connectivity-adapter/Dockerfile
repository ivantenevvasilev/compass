FROM --platform=$BUILDPLATFORM golang:1.21.1-alpine3.18 as builder

ENV BASE_APP_DIR /go/src/github.com/kyma-incubator/compass/components/connectivity-adapter
WORKDIR ${BASE_APP_DIR}

COPY go.mod go.sum ${BASE_APP_DIR}/
RUN go mod download -x

COPY . ${BASE_APP_DIR}

ARG TARGETOS TARGETARCH
RUN GOOS=$TARGETOS GOARCH=$TARGETARCH go build -v -o main ./cmd/main.go
RUN mkdir /app && mv ./main /app/main



FROM alpine:3.18.4
LABEL source = git@github.com:kyma-incubator/compass.git
WORKDIR /app

COPY --from=builder /app /app

CMD ["/app/main"]
