FROM golang:1.24-alpine AS builder
ENV CGO_ENABLED=0 \
    GO111MODULE=on \
    GOOS=linux
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o main ./main.go

FROM alpine:3.21
RUN adduser -D appuser
WORKDIR /app
COPY --from=builder /app/main .
USER appuser
EXPOSE 8080
ENTRYPOINT ["./main"]