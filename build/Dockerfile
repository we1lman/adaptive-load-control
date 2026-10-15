FROM golang:1.25-alpine AS build

WORKDIR /src
COPY go.mod ./
COPY internal ./internal
COPY cmd ./cmd
RUN CGO_ENABLED=0 go build -o /out/gateway ./cmd/gateway
RUN CGO_ENABLED=0 go build -o /out/loadgen ./cmd/loadgen

FROM alpine:3.20

COPY --from=build /out/gateway /bin/gateway
COPY --from=build /out/loadgen /bin/loadgen
