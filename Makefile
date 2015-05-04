PWD:=$(shell pwd)

build:
	docker build -t "gofirst" .
	docker run -v $(PWD):/opt/gofirst gofirst:latest
