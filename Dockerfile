# syntax=docker/dockerfile:1

ARG GO_VERSION=1.23.1
ARG TARGETOS TARGETARCH

# Use full Go image for build stage
FROM --platform=${BUILDPLATFORM} golang:${GO_VERSION} AS build

WORKDIR /fn
ENV CGO_ENABLED=0 GOFLAGS="-trimpath"

# Cache dependencies
RUN --mount=target=. \
    --mount=type=cache,target=/go/pkg/mod \
    go mod download

# Build with reproducibility flags
RUN --mount=target=. \
    --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    GOOS=${TARGETOS} GOARCH=${TARGETARCH} \
    go build -trimpath -ldflags="-s -w -buildid=" -o /function .

FROM gcr.io/distroless/base-debian12:nonroot@sha256:6ec5aa99dc335666e79dc64e4a6c8b89c33a543a1967f20d360922a80dd21f02 AS image

COPY --from=build --chmod=0755 /function /function

USER nonroot:nonroot
ENTRYPOINT ["/function"]
