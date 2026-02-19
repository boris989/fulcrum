# ----- build stage -----

FROM golang:1.25 AS build
WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -trimpath -ldflags="-s -w" \
    -o /out/fulcrum ./cmd/orders

# ----- runtime stage ------
FROM gcr.io/distroless/static:nonroot

WORKDIR /

COPY --from=build /out/fulcrum /fulcrum
COPY migrations /migrations
USER nonroot:nonroot

EXPOSE 8080

ENTRYPOINT ["/fulcrum"]