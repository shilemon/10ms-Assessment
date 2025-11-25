FROM golang:1.22 AS builder

WORKDIR /build
COPY app .
RUN go build -o server .

FROM alpine

COPY --from=builder /build/app/server /bin/server

EXPOSE 80
ENTRYPOINT ["/bin/server"]
