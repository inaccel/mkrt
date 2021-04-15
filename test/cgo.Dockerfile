# syntax=docker/dockerfile:1.2

FROM golang AS sysroot
WORKDIR /src/cgo
COPY cgo .
RUN go build -o a.out main.go

FROM inaccel/mkrt AS mkrt

ENV MKRT_SYSROOT_DIR=/host
ENV MKRT_CONFIG_PATH=/src

ENV MKRT_TOP_BUILD_DIR=/tmp

RUN --mount=from=sysroot,target=/host,ro mkrt

FROM scratch
COPY --from=mkrt /tmp /
ENTRYPOINT ["cgo/ld.so", "cgo/lib.so"]
