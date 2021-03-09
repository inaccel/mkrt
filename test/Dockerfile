# syntax=docker/dockerfile:1.2

FROM gcc AS sysroot
WORKDIR /src/hello-world
COPY hello-world .
RUN gcc main.c

FROM inaccel/mkrt AS mkrt

ENV MKRT_SYSROOT_DIR=/host
ENV MKRT_CONFIG_PATH=/src

ENV MKRT_TOP_BUILD_DIR=/tmp

RUN --mount=from=sysroot,target=/host,ro mkrt

FROM scratch
COPY --from=mkrt /tmp /
ENTRYPOINT ["hello-world/ld.so", "hello-world/lib.so"]
