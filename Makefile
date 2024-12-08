APP_NAME = editai
DOCKER_IMAGE = editai:latest
VERSION = 0.1.0
BUILD_TIME = $(shell date -u '+%Y-%m-%d_%H:%M:%S')
BUILD_OPTIONS := -clean -upx
BUILD_TAGS := 
GO_ENV := GOFLAGS="-buildvcs=false"

# Setup development environment
configure:
	docker compose up -d

dev:
	wails dev

# Build development version
build-dev:
	$(GO_ENV) wails build

# Build for different platforms
build-windows:
	$(GO_ENV) wails build -platform windows/amd64 -nsis ${BUILD_OPTIONS}

build-mac:
	$(GO_ENV) wails build -platform darwin/universal ${BUILD_OPTIONS}

build-linux:
	$(GO_ENV) wails build -platform linux/amd64 ${BUILD_OPTIONS} && \
		cd build/debian && \
		nfpm pkg --packager deb --target ../bin/ && \
		nfpm pkg --packager rpm --target ../bin/

# Clean up
clean:
	rm -rf build/bin/*
	docker compose down --rmi all

# Update dependencies
tidy:
	docker exec -e $(GO_ENV) editai-dev sh -c "go mod tidy && cd frontend && npm install"

.PHONY: configure build-dev build-windows build-mac build-linux build-linux-static clean tidy