# build
FROM golang:1.16-buster as builder

RUN apt-get update && apt-get install -y \
    libusb-1.0

RUN mkdir /build
WORKDIR /build
COPY . .

RUN go get -d -v ./...
RUN CGO_ENABLED=1 GOOS=linux go build -a -o kl-logger .

# run
FROM debian:buster

RUN apt-get update && apt-get install -y \
    libusb-1.0

COPY --from=builder /build/kl-logger .

EXPOSE 8088

ENTRYPOINT ["./kl-logger"]
CMD ["--apiPort 8088"]
