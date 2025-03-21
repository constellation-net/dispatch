FROM golang:1.23.2 AS builder
COPY . /app
WORKDIR /app
# See here for why this step is needed: https://github.com/docker-library/golang/issues/452
RUN git config --global --add safe.directory /app
RUN go build -o main .

# Multi-stage builds make the final Docker image much more space-efficient by removing unnecessary bloat (e.g: to Go compiler)
FROM alpine:3.21 AS runner
COPY --from=builder /app/main /app/main
ENTRYPOINT ["/app/main"]