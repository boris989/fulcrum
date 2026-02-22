# ----- build stage -----
FROM golang:1.25-bookworm AS build

RUN apt-get update && apt-get install -y \
    librdkafka-dev \
    gcc \
    g++

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .

ARG VERSION
ARG COMMIT
ARG BUILD_TIME

RUN go build \
    -ldflags "-s -w \
    -X github.com/boris989/fulcrum/internal/platform/version.Version=${VERSION} \
    -X github.com/boris989/fulcrum/internal/platform/version.Commit=${COMMIT} \
    -X github.com/boris989/fulcrum/internal/platform/version.BuildTime=${BUILD_TIME}" \
    -o /out/fulcrum ./cmd/orders

# ----- runtime stage ------
FROM gcr.io/distroless/cc-debian12

WORKDIR /

COPY --from=build /out/fulcrum /fulcrum
COPY migrations /migrations

EXPOSE 8080

ENTRYPOINT ["/fulcrum"]
