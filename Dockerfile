FROM golang:1.24-alpine AS builder

WORKDIR /app

COPY go.mod go.sum* ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /hrpolicy ./cmd/hrpolicy

FROM alpine:3.20

RUN apk add --no-cache ca-certificates

WORKDIR /app

COPY --from=builder /hrpolicy .
COPY documents/ ./documents/

ENV HR_DOCUMENTS_DIR=/app/documents

CMD ["./hrpolicy"]
