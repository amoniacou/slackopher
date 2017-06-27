BINARY=slackopher
VERSION=0.0.1
BUILD=`date +%FT%T%z`
LDFLAGS = -ldflags "-X main.Version=${VERSION} -X main.Build=${BUILD}"
NAME:=$(shell basename `git rev-parse --show-toplevel`)

all: build

build: clean
	pwd
	glide install
	cd ./cmd && go build ${LDFLAGS} -o ../bin/${BINARY}

clean:
	if [ -f ./bin/${BINARY} ] ; then rm ./bin/${BINARY} ; fi

docker-clean:
	docker rmi amoniacou/${NAME} &>/dev/null || true

docker-build: docker-clean
	docker build -t amoniacou/slackopher-build -f ./Dockerfile.build .
	docker run --rm -it -v "$$PWD":/go/src/github.com/amoniacou/slackopher -e CGO_ENABLED=true \
		-e LDFLAGS='-extldflags "-static"' \
		-e COMPRESS_BINARY=true \
		amoniacou/slackopher-build
	docker build --pull=true --no-cache -t amoniacou/${NAME} .

deps:
	glide install
