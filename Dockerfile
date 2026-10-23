# syntax=docker/dockerfile:1

FROM --platform=$BUILDPLATFORM golang:1.22-alpine AS build

ARG TARGETOS
ARG TARGETARCH

WORKDIR /src
COPY go.mod ./
COPY cmd/ ./cmd/
COPY internal/ ./internal/

RUN CGO_ENABLED=0 GOOS=${TARGETOS:-linux} GOARCH=${TARGETARCH:-amd64} \
    go build -trimpath -ldflags="-s -w" -o /out/aggregator ./cmd/aggregator

FROM alpine:3.22

RUN addgroup -S -g 10001 app && adduser -S -D -H -u 10001 -G app app
COPY --from=build /out/aggregator /aggregator

USER 10001:10001
ENTRYPOINT ["/aggregator"]
