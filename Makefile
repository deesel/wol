.PHONY: run all

.DEFAULT_GOAL := all

all:
	@docker buildx build -t deesel/wol:dev .

run:
	@docker run --rm -it --net=host -v $(shell pwd)/.dev:/etc/wol deesel/wol:dev
