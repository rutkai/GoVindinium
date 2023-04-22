FROM golang:1.20-buster as builder

WORKDIR /app

COPY go.* ./
RUN go mod download

COPY . ./
RUN go build -v -o vindinium-warrior

FROM debian:buster-slim

COPY --from=builder /app/vindinium-warrior /app/vindinium-warrior

CMD ["/app/vindinium-warrior"]
