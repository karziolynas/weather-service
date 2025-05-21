FROM golang:alpine AS builder
RUN apk update && apk add --no-cache git
RUN apk --no-cache add ca-certificates
WORKDIR $GOPATH/src/weather_service/
COPY . .

RUN go mod tidy
RUN go build -o /go/bin/weather_service

COPY ./certs /go/src/weather_service/certs

# FROM scratch
# COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
# COPY --from=builder /go/bin/weather_service /go/bin/weather_service
# COPY --from=builder /go/src/weather_service/index.html /go/src/weather_service/index.html
# COPY --from=builder /go/src/weather_service/certs /go/src/weather_service/certs
# ENTRYPOINT ["/go/bin/weather_service"]

FROM alpine
RUN apk --no-cache add ca-certificates
COPY --from=builder /go/bin/weather_service /go/bin/weather_service
COPY --from=builder /go/src/weather_service/certs /go/src/weather_service/certs
COPY --from=builder /go/src/weather_service/index.html /go/src/weather_service/index.html
ENTRYPOINT ["/go/bin/weather_service"]