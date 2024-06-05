FROM golang:1.22-alpine as builder
WORKDIR /code
COPY go.sum go.mod .
RUN go mod download -x

COPY *.go ./
RUN go build -v -o /livestream ./...

FROM alpine
COPY --from=builder /livestream /livestream
CMD ["/livestream"]
