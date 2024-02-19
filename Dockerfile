FROM golang:1.22.0-alpine as builder
RUN apk --update add ca-certificates
RUN cd ..
RUN mkdir lite-reader
WORKDIR lite-reader
COPY . ./
ENV GO111MODULE=on
RUN CGO_ENABLED=0 GOOS=linux go build -mod=vendor -a -installsuffix cgo -o lite-reader ./cmd/main.go

FROM scratch
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --from=builder /go/lite-reader/public ./public
COPY --from=builder /go/lite-reader/lite-reader .

# Directory to store the data, which can be referenced as the mounting point.
VOLUME /var/opt/lite-reader
EXPOSE 3000

ENTRYPOINT ["./lite-reader"]
