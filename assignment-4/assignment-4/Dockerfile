
FROM golang:1.21-alpine AS builder

WORKDIR /app

COPY . .

RUN go mod tidy && CGO_ENABLED=0 GOOS=linux go build -o myapp ./cmd/api/main.go


FROM alpine:3.19

WORKDIR /app

COPY --from=builder /app/myapp .

EXPOSE 8000

ENTRYPOINT ["./myapp"]