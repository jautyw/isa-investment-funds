FROM golang:1.23-alpine as build

ARG SERVICE='local'
ARG COMMIT='local'

WORKDIR /build

COPY . .

RUN apk --no-cache add ca-certificates && \
    GOOS=linux GOARCH=amd64 CGO_ENABLED=0 GO111MODULE=on go build -ldflags "-s -w -X main.service=${SERVICE} -X main.commit=${COMMIT}" -o /app ./cmd/server

FROM scratch

COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt

COPY --from=build /app /

COPY config-docker.yaml /

ENTRYPOINT ["/app"]