FROM golang:1.22.0 AS builder

WORKDIR /src

COPY . .

RUN go get -v . \
    && CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o license-checker

FROM alpine:3.19.1

WORKDIR /opt

COPY --from=builder /src/license-checker .

RUN chmod +x license-checker

CMD ["./license-checker"]
