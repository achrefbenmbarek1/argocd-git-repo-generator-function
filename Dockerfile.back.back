# syntax=docker/dockerfile:1

# ko requires scratch or distroless base
FROM gcr.io/distroless/static-debian12:nonroot
COPY --chmod=0755 ko-app/function /function
USER nonroot:nonroot
ENTRYPOINT ["/function"]
