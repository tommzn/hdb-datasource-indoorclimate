FROM --platform=linux/${TARGETARCH:-amd64} golang:1.19-alpine as builder

ARG TARGETPLATFORM
ARG BUILDPLATFORM
ARG TARGETOS
ARG TARGETARCH

ENV CGO_ENABLED=0
ENV GO111MODULE=on
ENV GOOS=${TARGETOS}
ENV GOARCH=${TARGETARCH}

WORKDIR /go/build

COPY go.mod go.mod
COPY go.sum go.sum
RUN go mod download

RUN apk -U add ca-certificates
RUN apk update && apk upgrade && apk add pkgconf git bash build-base sudo gcc musl-dev

COPY .  .

RUN go build -tags musl -ldflags '-extldflags "-static"' -o build_artifact_bin
RUN echo "${TARGETOS}/${TARGETARCH}" > build.arch

FROM busybox:1.35.0-uclibc as busybox

FROM --platform=linux/${TARGETARCH:-amd64} gcr.io/distroless/static:nonroot

ARG TARGETOS
ARG TARGETARCH

LABEL org.opencontainers.image.source=https://github.com/tommzn/hdb-datasource-indoorclimate

WORKDIR /go

COPY --from=busybox /bin/sh /bin/sh
COPY --from=busybox /bin/ls /bin/ls
COPY --from=busybox /bin/cat /bin/cat
COPY --from=busybox /bin/echo /bin/echo

COPY --from=builder /go/build/build_artifact_bin hdb-bin
COPY --from=builder /go/build/build.arch build.arch
RUN /bin/echo "${TARGETOS}/${TARGETARCH}" > runtime.arch
USER nonroot:nonroot

ENTRYPOINT ["/go/hdb-bin"]
