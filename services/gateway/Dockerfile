FROM golang:1.22-alpine AS builder

WORKDIR /app
COPY go.mod ./
COPY internal/ internal/
COPY cmd/ cmd/

RUN go build -o /app/gateway cmd/server/main.go

FROM alpine:latest
WORKDIR /root/
COPY --from=builder /app/gateway .

EXPOSE 8080
CMD ["./gateway"]
