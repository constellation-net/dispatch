FROM golang:1.23.2 AS builder
# See why --chown=default is needed here: https://groups.google.com/g/golang-nuts/c/LZbM2WlZoJM
COPY --chown=default . /app
WORKDIR /app
RUN go build -o main .

# Multi-stage builds make the final Docker image much more space-efficient by removing unnecessary bloat (e.g: to Go compiler)
FROM alpine:3.21 AS runner
COPY --from=builder /app/main /app/main
ENTRYPOINT ["/app/main"]