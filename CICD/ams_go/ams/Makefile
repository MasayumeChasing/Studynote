app=ams
docker_hub=registry-cbu.huawei.com/holo-fuxi
CID_REPO_COMMIT=$(shell git rev-parse --short HEAD)
artifactTag=$(shell date '+%Y.%m%d')
TAG=$(artifactTag)_$(CID_GLOBAL_REPO_BRANCH)_$(CID_REPO_COMMIT)       

version:
	cat .build_config/imageVersion.text

build: CGO_ENABLED=1
build: CGO_CFLAGS="-fstack-protector-all -ftrapv -D_FORTIFY_SOURCE=2 -O2 -Wall"
build: CGO_CPPFLAGS="-fstack-protector-all -ftrapv -D_FORTIFY_SOURCE=2 -O2 -Wall"
build: CGO_LDFLAGS="-Wl,-s,-z,now,--build-id=none"
build:
	go build -buildmode=default -ldflags="-w -s -X main.Version=1.0.0 -X 'main.BuildTime=`date`' -X 'main.GoVersion=`go version`'" -o dockerfile/$(app) -v 

image: build
	mkdir -p .build_config
	echo "$(docker_hub)/$(app):$(TAG)" > .build_config/imageVersion.text
	chmod 500 dockerfile/$(app)
	chmod 500 dockerfile/entrypoint.sh
	docker build -t $(docker_hub)/$(app):$(TAG) dockerfile

clean:
	rm -rf .build_config

.PHONY: build
