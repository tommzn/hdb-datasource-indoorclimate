FROM alpine:latest

RUN apk add --no-cache ca-certificates

WORKDIR /go

COPY --chmod=0755 build_artifact_bin k8s-bin

RUN chmod 755 /go/k8s-bin
ENTRYPOINT ["/go/k8s-bin"]