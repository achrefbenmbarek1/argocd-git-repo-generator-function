# syntax=docker/dockerfile:1

ARG GO_VERSION=1.23.1
ARG TARGETOS
ARG TARGETARCH

# Use full Go image for build stage
FROM --platform=${BUILDPLATFORM} golang:${GO_VERSION} AS build

WORKDIR /fn
ENV CGO_ENABLED=0 

# Cache dependencies
RUN --mount=target=. \
    --mount=type=cache,target=/go/pkg/mod \
    go mod download

# Build with reproducibility flags
RUN --mount=target=. \
    --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    GOOS=${TARGETOS} GOARCH=${TARGETARCH} \
    go build -o /function .

FROM gcr.io/distroless/static-debian12:nonroot@sha256:6ec5aa99dc335666e79dc64e4a6c8b89c33a543a1967f20d360922a80dd21f02 AS image

WORKDIR /
COPY --from=build --chmod=0755 /function /function
EXPOSE 9443
USER nonroot:nonroot
ENTRYPOINT ["/function"]
