FROM golang:alpine AS builder
WORKDIR /go-hello
COPY go.mod .
COPY go.sum .
RUN go mod download
COPY src/ src
RUN go build -o server src/main.go

FROM alpine
WORKDIR /
EXPOSE 8080/tcp
RUN mkdir /views
COPY views /views
COPY --from=builder /go-hello/server .
CMD ["/server"]