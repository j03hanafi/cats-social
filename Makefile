# Variables
MAIN_PACKAGE_PATH := ./cmd/api
BINARY_NAME := trainer-service-backend


# ==================================================================================== #
# HELPERS
# ==================================================================================== #

## help: print this help message
.PHONY: help
help:
	@echo 'Usage:'
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' |  sed -e 's/^/ /'


# ==================================================================================== #
# QUALITY CONTROL
# ==================================================================================== #

## tidy: run formatting, go mod tidy and goimports
.PHONY: tidy
tidy:
	go run mvdan.cc/gofumpt@latest -extra -l -w .
	go run github.com/segmentio/golines@latest --max-len=120 --shorten-comments -w .
	go run github.com/incu6us/goimports-reviser/v3@latest -rm-unused ./...
	go mod tidy -v

## audit: run quality control checks
.PHONY: audit
audit:
	go mod verify
	go vet ./...
	go run github.com/roblaszczak/go-cleanarch@latest -application service -interfaces handler -infrastructure repository
	go run go.uber.org/nilaway/cmd/nilaway@latest ./...
	go run github.com/golangci/golangci-lint/cmd/golangci-lint@latest run
	go test -race -buildvcs -vet=off ./...


# ==================================================================================== #
# DEVELOPMENT
# ==================================================================================== #

## test: run all tests
.PHONY: test
test:
	go test -v -race -buildvcs -count=1 ./...

## test/cover: run all tests and display coverage
.PHONY: test/cover
test/cover:
	go test -v -race -buildvcs -count=1 -coverprofile=./tmp/coverage.out ./...
	go tool cover -html=./tmp/coverage.out

## build: build the application
.PHONY: build
build:
	# Include additional build steps, like TypeScript, SCSS or Tailwind compilation here...
	go build -o=./tmp/bin/${BINARY_NAME} ${MAIN_PACKAGE_PATH}

## run: run the  application
.PHONY: run
run: build
	./tmp/bin/${BINARY_NAME}

## watch: run the application with reloading on file changes
.PHONY: watch
watch:
	go run github.com/cosmtrek/air@latest \
		--build.cmd "make build" --build.bin "./tmp/bin/${BINARY_NAME}" --build.delay "100" \
		--build.exclude_dir "" \
		--build.include_ext "go, mod, tpl, tmpl, html, env" \
		--build.send_interrupt "true" \
		--build.kill_delay "5000000" \
		--misc.clean_on_exit "true"

## clean: remove the binary
.PHONY: clean
clean:
	rm -f ./tmp/bin/${BINARY_NAME}