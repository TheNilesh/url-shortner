FROM golang:1.20 AS builder
WORKDIR /app
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o url-shortener

FROM alpine:latest
RUN adduser -D -g '' nonroot
WORKDIR /home/nonroot
COPY --from=builder /app/url-shortener /usr/local/bin/url-shortener
COPY config.properties .

EXPOSE 8080
USER nonroot:nonroot
CMD ["url-shortener"]
