# Builder stage

FROM golang:1.26-alpine AS builder

RUN apk add --no-cache gcc musl-dev

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN ./setup.sh

RUN CGO_ENABLED=1 GOOS=linux go build -o seaurl ./cmd/api/main.go

# Runner stage

FROM alpine:latest

RUN apk --no-cache add ca-certificates tzdata

WORKDIR /app

COPY --from=builder /app/seaurl .

COPY --from=builder /app/migrations ./migrations

EXPOSE 2900

CMD ["./seaurl"]