app=mps
REPOSITORY=registry-cbu.huawei.com
ORG ?= ivm-dev
CID_REPO_COMMIT=$(shell git rev-parse --short HEAD)
artifactTag=$(shell date '+%Y.%m%d')
TAG=$(artifactTag)_$(CID_GLOBAL_REPO_BRANCH)_$(CID_REPO_COMMIT)

version:
	cat .build_config/imageVersion.text

build:
	cd deps && sh ./gen_deps.sh && cd ..
	cmake .
	make -j4
	strip -s bin/holo_$(app)
	mv bin/holo_$(app) dockerfile/holo_$(app)
	mv bin/libmps_common.so dockerfile/libmps_common.so
	cp bin/get_version.sh dockerfile/get_version.sh
	echo -e "build_time:`date --rfc-3339=seconds`\nbuild_branch:${CID_GLOBAL_REPO_BRANCH}\nbuild_commit:`git rev-parse HEAD`\nbuild_user:${CID_BUILD_USER}" > dockerfile/build.version
image: build
	mkdir -p .build_config
	echo "$(REPOSITORY)/$(ORG)/$(app):$(TAG)" > .build_config/imageVersion.text
	chmod 500 dockerfile/holo_$(app)
	chmod 500 dockerfile/entrypoint.sh
	chmod 500 dockerfile/holo_amah_agent
	chmod 500 dockerfile/get_version.sh
	chmod 500 dockerfile/build.version
	docker build -t $(REPOSITORY)/$(ORG)/$(app):$(TAG) dockerfile

clean:
	rm -rf .build_config

.PHONY: build
