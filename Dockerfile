FROM --platform=$BUILDPLATFORM golang:1.23-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

ARG TARGETOS
ARG TARGETARCH
RUN CGO_ENABLED=0 GOOS=$TARGETOS GOARCH=$TARGETARCH go build -o volume_toolkit .

FROM alpine

RUN apk add --no-cache curl

COPY --from=builder /app/volume_toolkit /usr/bin/volume_toolkit
RUN chmod +x /usr/bin/volume_toolkit

ENTRYPOINT ["/usr/bin/volume_toolkit"]
