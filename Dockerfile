FROM golang:alpine as builder
WORKDIR /go/src/github.com/atotto/autocert
RUN apk add --no-cache git
RUN go get -d golang.org/x/crypto/acme/autocert 
COPY main.go .
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags '-w -s' -a -installsuffix cgo -o app .

FROM alpine:latest
RUN apk add --no-cache ca-certificates
COPY --from=builder /go/src/github.com/atotto/autocert/app /
EXPOSE 443
EXPOSE 80
ENTRYPOINT ["/app"]
