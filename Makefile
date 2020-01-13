PACKAGE_VERSION ?= latest
REGISTRY ?= dockerhub.com/alleeclark/git-events
ARGS ?= start
LIBGITVERSION = v0.28.1

.PHONY: images
images:
	docker build -t git-events:$(PACKAGE_VERSION) .

.PHONY: runserver
runserver:
	docker run -it git-events:$(PACKAGE_VERSION) start