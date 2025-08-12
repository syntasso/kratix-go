FROM --platform=$BUILDPLATFORM golang:1.24 AS builder
ARG TARGETARCH
ARG TARGETOS

WORKDIR /workspace

ENV GO111MODULE=on
ENV CGO_ENABLED=0
ENV GOOS=${TARGETOS}
ENV GOARCH=${TARGETARCH}

# Build
RUN --mount=target=. \
    --mount=type=cache,target=/root/.cache/go-build \
    --mount=type=cache,target=/go/pkg \
    go build -a -o /out/resource-pipeline system/assets/workflows/resource/main.go
RUN --mount=target=. \
    --mount=type=cache,target=/root/.cache/go-build \
    --mount=type=cache,target=/go/pkg \
    go build -a -o /out/promise-pipeline system/assets/workflows/promise/main.go

FROM gcr.io/distroless/cc:nonroot
WORKDIR /
COPY --from=builder /out/resource-pipeline .
COPY --from=builder /out/promise-pipeline .
USER 65532:65532

ENTRYPOINT []
