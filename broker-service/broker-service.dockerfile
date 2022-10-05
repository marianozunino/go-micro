FROM golang:alpine AS builder

RUN mkdir /app
COPY . /app
WORKDIR /app
RUN CGO_ENABLED=0 GOOS=linux go build -o brokerApp ./cmd/api

RUN chmod +x brokerApp

# tiny image

FROM alpine:latest

RUN mkdir /app
COPY --from=builder /app/brokerApp /app

CMD ["/app/brokerApp"]



