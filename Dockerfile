FROM golang:alpine AS builder
WORKDIR /root
ADD . .
WORKDIR /root/cmd/passport
RUN env CGO_ENABLED=0 go build -v -trimpath -ldflags '-s -w'
FROM scratch
COPY --from=builder /root/cmd/passport/passport /passport
ENTRYPOINT ["/passport"]
