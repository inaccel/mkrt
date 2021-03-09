FROM golang:alpine AS rtldd
WORKDIR /go/src/github.com/inaccel/mkrt/rtldd
COPY rtldd/go.mod .
COPY rtldd/go.sum .
RUN go mod download
COPY rtldd .
RUN go build -o /go/bin/rtldd
