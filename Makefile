.PHONY: test build
SHELL   := bash
NAME    := navikt/vks
LATEST  := ${NAME}:latest

push-dockerhub: docker-build
docker-build: build-linux
build: test

clean:
	rm vks

test:
	go test ./...
build:
	CGO_ENABLED=0 go build -a -installsuffix cgo -o vks
build-linux:
	CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o vks

docker-build:
	docker image build -t ${NAME}:$(CIRCLE_BUILD_NUM) -t ${LATEST} -f Dockerfile .

push-dockerhub:
	docker image push ${NAME}:$(CIRCLE_BUILD_NUM) 
