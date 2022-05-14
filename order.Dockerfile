FROM golang:1.18-alpine AS builder
LABEL stage=gobuilder

ENV CGO_ENABLED 0
ENV GOOS linux
RUN apk update --no-cache && apk add --no-cache tzdata
WORKDIR /build
ADD go.mod .
ADD go.sum .
RUN go mod download
COPY ../.. .
RUN go build -ldflags="-s -w" -o /app/order ./cmd/store-order/main.go

FROM alpine
RUN apk update --no-cache && apk add --no-cache ca-certificates

WORKDIR /app

COPY --from=builder /app/order /app/order
RUN mkdir configs
COPY --from=builder /build/configs /app/configs
CMD ["./order", "-c", "configs/order.yaml"]