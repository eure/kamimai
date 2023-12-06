FROM golang:1.21-buster as builder
WORKDIR /go
COPY ./ .
RUN export GOPATH= && go mod download && go build -o /go/ -v -ldflags '-s -w' ./cmd/kamimai

FROM gcr.io/distroless/base
COPY --from=builder go/kamimai .
ENTRYPOINT ["/kamimai"]
CMD ["--help"]
