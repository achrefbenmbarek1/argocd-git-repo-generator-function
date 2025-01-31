# syntax=docker/dockerfile:1.4

# Pin exact Go version and digest for reproducibility
ARG GO_VERSION=1.22.4
ARG GO_DIGEST=sha256:0f5f3a28f710a5bbab5b31a7f8b6463fe14a1a8d0656f3e7e0b5a15dcf6c7f1e

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
FROM gcr.io/distroless/static-debian12@sha256:ad2fc03f2c995491b6d2eb357fe1a6d05d47c1f8b1160a1a38641b3c405a7f6e

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
