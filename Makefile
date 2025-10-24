# Change these variables as necessary.
main_package_path = ./
binary_name = gentmpl
binary_fullpath = ${binary_name}

# ==================================================================================== #
# Build Version Information
# ==================================================================================== #

build build-prod: verinfopath = github.com/mmbros/gentmpl/internal/version
build build-prod: GO_VERSION := $(shell go version)
build build-prod: GIT_COMMIT := $(shell git rev-parse --short HEAD)
build build-prod: BUILD_TIME := $(shell date '+%F %T %z')
build build-prod: OS_ARCH := $(shell uname -s -m)

build-prod: APP_VERSION_PROD := $(shell git tag | grep ^v | sort -V | tail -n 1)
build: APP_VERSION_DEV := dev-$(shell date +%Y%m%dT%H%M%S)

build build-prod: VERSION_INFO_COMMON = -X '${verinfopath}.BuildTime=${BUILD_TIME}' \
                        -X '${verinfopath}.GitCommit=${GIT_COMMIT}' \
				        -X '${verinfopath}.GoVersion=${GO_VERSION}' \
				        -X '${verinfopath}.OsArch=${OS_ARCH}'


build : VERSION_INFO_DEV = -X '${verinfopath}.AppVersion=${APP_VERSION_DEV}' ${VERSION_INFO_COMMON}
						
build-prod : VERSION_INFO_PROD = -X '${verinfopath}.AppVersion=${APP_VERSION_PROD}' ${VERSION_INFO_COMMON}


# ==================================================================================== #
# HELPERS
# ==================================================================================== #

## help: print this help message
.PHONY: help
help:
	@echo 'Usage:'
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' |  sed -e 's/^/ /'

.PHONY: confirm
confirm:
	@echo -n 'Are you sure? [y/N] ' && read ans && [ $${ans:-N} = y ]

.PHONY: no-dirty
no-dirty:
	@test -z "$(git status --porcelain)"

# ==================================================================================== #
# QUALITY CONTROL
# ==================================================================================== #

## audit: run quality control checks
.PHONY: audit
audit: test
	go mod tidy -diff
	go mod verify
	test -z "$(gofmt -l .)" 
	go vet ./...
	go run honnef.co/go/tools/cmd/staticcheck@latest -checks=all,-ST1000,-U1000 ./...
	go run golang.org/x/vuln/cmd/govulncheck@latest ./...

## test: run all tests
.PHONY: test
test:
	go test -v -race -buildvcs ./...

## test/cover: run all tests and display coverage
.PHONY: test/cover
test/cover:
	go test -v -race -buildvcs -coverprofile=/tmp/coverage.out ./...
	go tool cover -html=/tmp/coverage.out

## upgradeable: list direct dependencies that have upgrades available
.PHONY: upgradeable
upgradeable:
	@go run github.com/oligot/go-mod-upgrade@latest

# ==================================================================================== #
# DEVELOPMENT
# ==================================================================================== #

## tidy: tidy modfiles and format .go files
.PHONY: tidy
tidy:
	go mod tidy -v
	go fmt ./...

## build: build the application
.PHONY: build
build:
	# Include additional build steps, like TypeScript, SCSS or Tailwind compilation here...
	go build -ldflags "${VERSION_INFO_DEV}" -o=${binary_fullpath} ${main_package_path}

.PHONY: build-prod
build-prod:
	# Include additional build steps, like TypeScript, SCSS or Tailwind compilation here...
	go build -ldflags "-s ${VERSION_INFO_PROD}" -o=${binary_fullpath} ${main_package_path}


## run: run the  application
.PHONY: run
run: build
	${binary_fullpath}

## run/live: run the application with reloading on file changes
.PHONY: run/live
run/live:
	go run github.com/cosmtrek/air@v1.43.0 \
	    --build.cmd "make build" --build.bin "${binary_fullpath}" --build.delay "100" \
	    --build.exclude_dir "" \
	    --build.include_ext "go, tpl, tmpl, html, css, scss, js, ts, sql, jpeg, jpg, gif, png, bmp, svg, webp, ico" \
	    --misc.clean_on_exit "true"

# ==================================================================================== #
# OPERATIONS
# ==================================================================================== #

## push: push changes to the remote Git repository
.PHONY: push
push: confirm audit no-dirty
	git push

## production/deploy: deploy the application to production
.PHONY: production/deploy
production/deploy: confirm audit no-dirty build-prod
# 	GOOS=linux GOARCH=amd64 go build -ldflags="-s ${VERSION_INFO_PROD}" -o=/tmp/bin/linux_amd64/${binary_name} ${main_package_path}
	upx -5 ${binary_fullpath}
	# Include additional deployment steps here...

