# build stage
FROM golang:1.20-alpine AS builder
RUN apk add --no-cache ca-certificates git build-base
WORKDIR /src
COPY ./app ./
ENV CGO_ENABLED=0
RUN go mod download
RUN go build -o /out/app ./server.go

# runtime stage
FROM alpine:3.18
RUN apk add --no-cache ca-certificates
# create non-root user
RUN addgroup -S appgroup && adduser -S appuser -G appgroup
WORKDIR /app
COPY --from=builder /out/app /app/app
# ensure correct port
EXPOSE 8080
USER appuser
ENTRYPOINT ["/app/app"]
