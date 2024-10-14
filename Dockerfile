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
RUN apk add --no-cache tzdata mpg123
COPY assets/siren.mp3 /app/assets/siren.mp3
COPY --from=builder /build/backend /app/backend
WORKDIR /app
CMD ["./backend"]

# RUN GORACE="halt_on_error=1" CGO_ENABLED=1 GOOS=linux go build -race -o /build/backend ./cmd/backend