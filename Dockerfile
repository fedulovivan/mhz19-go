FROM golang:alpine AS builder
RUN apk --no-cache add gcc musl-dev
WORKDIR /build
COPY go.mod go.sum ./
RUN go mod download
COPY cmd cmd
COPY pkg pkg
COPY internal internal
RUN CGO_ENABLED=1 GOOS=linux go build -o /build/backend ./cmd/backend

FROM alpine:latest
COPY --from=builder /build/backend /backend
CMD ["/backend"]