FROM golang:1.17.10 as builder

RUN mkdir -p /shortener
COPY . /shortener
WORKDIR /shortener

RUN GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o api ./cmd

FROM scratch

COPY --from=builder /shortener /shortener
COPY --from=builder /etc/ssl/certs /etc/ssl/certs

ENTRYPOINT ["/shortener/api"]