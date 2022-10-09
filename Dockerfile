FROM golang:1.19-alpine AS build

RUN apk add --no-cache git

WORKDIR /build

COPY go.* ./
RUN go mod download

COPY . ./
RUN CGO_ENABLED=0 go build -ldflags="-s -w" -o /app ./cmd/api

FROM alpine:latest

RUN apk add --no-cache ca-certificates
COPY --from=build /app /app
COPY --from=build /build/migrations /migrations
ENTRYPOINT ["/app"]
