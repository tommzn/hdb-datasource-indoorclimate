FROM alpine:latest

RUN apk add --no-cache ca-certificates

WORKDIR /go

COPY build_artifact_bin k8s_bin

RUN chmod 755 /go/k8s_bin
ENTRYPOINT ["/go/k8s_bin"]
