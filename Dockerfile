FROM golang:1.24-alpine AS builder

ENV GOTOOLCHAIN=auto

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY *.go ./
RUN CGO_ENABLED=0 GOOS=linux go build -o server .

FROM alpine:3.20

COPY --from=builder /app/server /server

EXPOSE 8080

CMD ["/server"]
