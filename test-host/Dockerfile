FROM golang:bullseye AS test-host-builder
COPY . /services
WORKDIR /services
RUN go mod tidy && CGO_ENABLED=0 go build -o api

FROM alpine:latest AS test-host-run
COPY --from=test-host-builder /services/api /
USER root
RUN ( getent group deploy || addgroup -S deploy ) && \
  ( getent passwd deploy || adduser -S deploy -G deploy ) && \
  chown -R deploy:deploy /api

FROM test-host-run
USER deploy
EXPOSE 9100
CMD ["/api"]
