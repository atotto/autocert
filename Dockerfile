FROM golang:alpine as builder
WORKDIR /workspace
ADD . /workspace
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags '-w -s' -a -installsuffix cgo -o app .

FROM alpine:latest
RUN apk add --no-cache ca-certificates
COPY --from=builder /workspace/app /
EXPOSE 443
EXPOSE 80
ENTRYPOINT ["/app"]
