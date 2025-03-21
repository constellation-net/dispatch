FROM golang:1.23.2 AS builder
USER root
COPY . /app
WORKDIR /app
RUN chown -R 0:0 /app
RUN go build -v -o main .

# Multi-stage builds make the final Docker image much more space-efficient by removing unnecessary bloat (e.g: to Go compiler)
FROM alpine:3.21 AS runner
COPY --from=builder /app/main /app/main
ENTRYPOINT ["/app/main"]