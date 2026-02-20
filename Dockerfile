# ----- build stage -----
FROM golang:1.25 AS build
WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN GOOS=linux GOARCH=amd64 \
    go build -trimpath -ldflags="-s -w" \
    -o /out/fulcrum ./cmd/orders

# ----- runtime stage ------
FROM gcr.io/distroless/base-debian12

WORKDIR /

COPY --from=build /out/fulcrum /fulcrum
COPY migrations /migrations
USER nonroot:nonroot

EXPOSE 8080

ENTRYPOINT ["/fulcrum"]
