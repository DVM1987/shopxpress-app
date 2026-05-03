# syntax=docker/dockerfile:1.7

# ============================================================
# Stage 1: builder — compile Go binary
# ============================================================
FROM --platform=$BUILDPLATFORM golang:1.26-alpine AS builder

ARG TARGETOS
ARG TARGETARCH
ARG SERVICE=gateway
ARG GIT_COMMIT=unknown
ARG IMAGE_TAG=dev

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY pkg/ ./pkg/
COPY services/${SERVICE}/ ./services/${SERVICE}/

RUN CGO_ENABLED=0 GOOS=${TARGETOS} GOARCH=${TARGETARCH} \
    go build \
      -trimpath \
      -ldflags="-s -w -X main.imageTag=${IMAGE_TAG} -X main.gitCommit=${GIT_COMMIT}" \
      -o /out/app \
      ./services/${SERVICE}

# ============================================================
# Stage 2: final — distroless, non-root, runtime only
# ============================================================
FROM gcr.io/distroless/static-debian12:nonroot

ARG SERVICE=gateway
ARG IMAGE_TAG=dev
ARG GIT_COMMIT=unknown

LABEL org.opencontainers.image.title="shopxpress-${SERVICE}"
LABEL org.opencontainers.image.source="https://github.com/dauvanmuoi/shopxpress-app"
LABEL org.opencontainers.image.revision="${GIT_COMMIT}"
LABEL org.opencontainers.image.version="${IMAGE_TAG}"

WORKDIR /

COPY --from=builder /out/app /app

ENV PORT=8080 \
    IMAGE_TAG=${IMAGE_TAG} \
    GIT_COMMIT=${GIT_COMMIT}

USER 65532:65532

EXPOSE 8080

ENTRYPOINT ["/app"]
