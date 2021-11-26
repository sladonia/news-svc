FROM golang:1.17-alpine AS builder
RUN apk --no-cache add git
WORKDIR /build
COPY go.mod go.sum ./
RUN go env -w GOPROXY=direct && go env -w GOSUMDB=off
RUN mkdir ~/.ssh && echo "StrictHostKeyChecking no" >> ~/.ssh/config
RUN go mod download

COPY cmd ./cmd
COPY internal ./internal
RUN CGO_ENABLED=0 GOOS=linux go build -o ./bin/serve -ldflags="-s -w" ./cmd/serve/*.go

FROM alpine:3.14.2
WORKDIR /app
COPY --from=builder /build/bin/serve .
CMD ["./serve"]
