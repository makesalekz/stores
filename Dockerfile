FROM golang:1.23 AS builder

WORKDIR /src
COPY . .
RUN go mod download
RUN CGO_ENABLED=0 go build -ldflags "-s -w" -o /app ./cmd/app/

FROM alpine:3.19
RUN apk --no-cache add ca-certificates
COPY --from=builder /app /app
COPY configs/config.local.yaml /config.yaml

EXPOSE 8000 9000
ENTRYPOINT ["/app"]
CMD ["-conf", "/config.yaml"]
