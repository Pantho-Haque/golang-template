FROM golang:1.25-alpine AS builder

WORKDIR /app

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o /prism ./cmd/workers/main.go

# Final stage
FROM alpine:latest

RUN apk add --no-cache ca-certificates tzdata

RUN adduser -D -g '' appuser

WORKDIR /app

COPY --from=builder /prism .

USER appuser

EXPOSE 80

CMD ["./prism"] 
