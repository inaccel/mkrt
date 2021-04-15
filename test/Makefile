all: cgo hello-world

cgo: build-cgo run-cgo

hello-world: build-hello-world run-hello-world

build-%:
	DOCKER_BUILDKIT=1 docker build --file $*.Dockerfile --tag $* .

run-%:
	docker run --interactive --rm --tty $*
