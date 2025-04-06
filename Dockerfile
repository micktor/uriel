FROM golang:1.23-alpine AS builder
RUN mkdir /app
ADD . /app
WORKDIR /app
RUN go build -o main .

FROM alpine
WORKDIR /app
COPY --from=builder /app/main /app/main

CMD ["/app/main", "httpd"]