# GO Dockerfile build

FROM golang:1.19-alpine AS builder

WORKDIR /go/src/app

COPY . .

RUN go get -d -v ./...

RUN go build .

FROM alpine

WORKDIR /app

COPY --from=builder /go/src/app/cloudflare_bgp_announcement_exporter .

EXPOSE 8080

CMD ["./cloudflare_bgp_announcement_exporter"]
