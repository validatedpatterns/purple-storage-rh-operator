FROM registry.access.redhat.com/ubi9/ubi-minimal:latest AS builder
ARG TARGETOS
ARG TARGETARCH
RUN microdnf install git-core jq tar -y && microdnf clean all

# Build the manager binary

WORKDIR /workspace

# Copy the go source
COPY go.mod go.mod
COPY go.sum go.sum

# use latest Go z release
ENV GOTOOLCHAIN=auto
ENV GO_INSTALL_DIR=/golang

# Ensure correct Go version
RUN export GO_VERSION=$(grep -E "go [[:digit:]]\.[[:digit:]][[:digit:]]" go.mod | awk '{print $2}') && \
    export GO_FILENAME=$(curl -sL 'https://go.dev/dl/?mode=json&include=all' | jq -r "[.[] | select(.version | startswith(\"go${GO_VERSION}\"))][0].files[] | select(.os == \"linux\" and .arch == \"amd64\") | .filename") && \
    curl -sL -o go.tar.gz "https://golang.org/dl/${GO_FILENAME}" && \
    mkdir -p ${GO_INSTALL_DIR} && \
    tar -C ${GO_INSTALL_DIR} -xzf go.tar.gz && \
    rm go.tar.gz && ln -sf ${GO_INSTALL_DIR}/go/bin/go /usr/bin/go

# cache deps before building and copying source so that we don't need to re-download as much
# and so that source changes don't invalidate our downloaded layer
RUN go mod download

# Copy the go source
COPY vendor/ vendor/
COPY cmd/main.go cmd/main.go
COPY api/ api/
COPY internal/ internal/
COPY files/ /files/
COPY hack/ hack/
# Needed to get the git versions in there
COPY .git/ .git/
RUN mkdir /licenses
COPY LICENSE /licenses

# Build
RUN --mount=type=secret,id=pull hack/build.sh

# UBI is larger (158Mb vs. 56Mb) but approved by RH 
FROM registry.access.redhat.com/ubi9/ubi-minimal:latest
WORKDIR /
COPY --from=builder /files/ /files/
COPY --from=builder /workspace/manager .
COPY --from=builder /licenses/ /licenses/
USER 65532:65532

ENTRYPOINT ["/manager"]
