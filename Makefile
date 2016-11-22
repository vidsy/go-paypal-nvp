BRANCH = "master"
VERSION = $(shell cat ./VERSION)
GO_BUILDER_IMAGE ?= "vidsyhq/go-builder"
PATH_BASE ?= "/go/src/github.com/vidsy"
REPONAME ?= "go-paypalnvp"

DEFAULT: test

check-version:
	git fetch
	(! git rev-list ${VERSION})

push-tag:
	git checkout ${BRANCH}
	git pull origin ${BRANCH}
	git tag ${VERSION}
	git push origin ${BRANCH} --tags

test:
	@docker run \
	-it \
	--rm \
	-v "${CURDIR}":${PATH_BASE}/${REPONAME} \
	-w ${PATH_BASE}/${REPONAME} \
	--entrypoint=go \
	${GO_BUILDER_IMAGE} test . ./nvppayload

test_ci:
	@docker run \
	-v "${CURDIR}":${PATH_BASE}/${REPONAME} \
	-w ${PATH_BASE}/${REPONAME} \
	--entrypoint=go \
	${GO_BUILDER_IMAGE} test . ./nvppayload -cover
