# syntax=docker/dockerfile:1.4

# Pin exact Go version and digest for reproducibility
ARG GO_VERSION=1.22.4
ARG GO_DIGEST=sha256:c8736b8dbf2b12c98bb0eeed91eef58ecef52b8c2bd49b8044531e8d8d8d58e8

# Use explicit platform and digest for build stage
FROM --platform=${BUILDPLATFORM} golang:${GO_VERSION}@${GO_DIGEST} AS build

WORKDIR /fn

# Reproducibility configurations
ENV CGO_ENABLED=0 \
    GOFLAGS="-trimpath -mod=readonly" \
    TZ=UTC

# Create fixed timestamp for build
ARG SOURCE_DATE_EPOCH
RUN <<EOF
    export DEBIAN_FRONTEND=noninteractive
    apt-get update && apt-get install -y --no-install-recommends \
    $(sort -u packages.txt) && \
    rm -rf /var/lib/apt/lists/*
    find /go/pkg/mod -exec touch -d @${SOURCE_DATE_EPOCH} {} +
EOF

# Cache dependencies with sorted modules
RUN --mount=target=. \
    --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    go mod download -x all

# Build with reproducibility flags
ARG TARGETOS
ARG TARGETARCH
RUN --mount=target=. \
    --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    GOOS=${TARGETOS} GOARCH=${TARGETARCH} \
    go build -trimpath -ldflags="-buildid= -w -s" -o /function .

# Final image with pinned distroless digest
FROM gcr.io/distroless/static-debian12@sha256:3f2b64ef97bd285e36132c684e6b2ae8f2723293d09aae046196cca64251acac

WORKDIR /
COPY --from=build --chmod=0755 /function /function

# Set fixed timestamps for all files
ARG SOURCE_DATE_EPOCH
RUN <<EOF
    find /function -exec touch -d @${SOURCE_DATE_EPOCH} {} +
    chown -R nonroot:nonroot /function
EOF

EXPOSE 9443
USER nonroot:nonroot
ENTRYPOINT ["/function"]
