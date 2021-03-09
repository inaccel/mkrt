FROM golang:alpine AS rtldd
WORKDIR /go/src/github.com/inaccel/mkrt/rtldd
COPY rtldd/go.mod .
COPY rtldd/go.sum .
RUN go mod download
COPY rtldd .
RUN go build -o /go/bin/rtldd

FROM alpine
RUN apk add --no-cache \
	binutils \
	coreutils \
	findutils \
	grep \
	patchelf \
	pkgconfig
COPY --from=rtldd /go/bin/rtldd /bin/rtldd
COPY mkrt /bin/mkrt
ENTRYPOINT ["mkrt"]
