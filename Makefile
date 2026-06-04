.PHONY: help build-backend build-frontend sync-frontend build start dev clean docker

SHELL := /bin/bash

VERSION ?= latest
IMG_NAME ?= fuming/standby_web

BACKEND_OUT ?= bin/server
WEB_DIR := web
WEB_DIST := $(WEB_DIR)/dist
WEB_ASSETS_DIR := assets/web/v1

help:
	@echo "Targets:"
	@echo "  build-backend     Build Go server -> $(BACKEND_OUT)"
	@echo "  build-frontend    Build Vue frontend -> $(WEB_DIST)"
	@echo "  sync-frontend     Copy $(WEB_DIST) -> $(WEB_ASSETS_DIR)"
	@echo "  build             Build backend+frontend and sync assets"
	@echo "  start             Build+sync then run server (http://127.0.0.1:1423/v1/web/)"
	@echo "  dev               Run server (go run) + Vite dev server"
	@echo "  clean             Remove build artifacts"
	@echo "  docker            Build image $(IMG_NAME):$(VERSION)"

build-backend:
	@set -euo pipefail; \
	mkdir -p "$$(dirname "$(BACKEND_OUT)")"; \
	go build -o "$(BACKEND_OUT)" .

build-frontend:
	@set -euo pipefail; \
	cd "$(WEB_DIR)"; \
	npm install; \
	npm run build

sync-frontend:
	@set -euo pipefail; \
	if [ ! -d "$(WEB_DIST)" ]; then \
		echo "Missing $(WEB_DIST). Run: make build-frontend"; \
		exit 1; \
	fi; \
	mkdir -p "$(WEB_ASSETS_DIR)"; \
	rm -rf "$(WEB_ASSETS_DIR)"/*; \
	cp -R "$(WEB_DIST)"/. "$(WEB_ASSETS_DIR)"/

build: build-backend build-frontend sync-frontend

start: build
	@set -euo pipefail; \
	if [ -f .env ]; then set -a && . ./.env && set +a; fi; \
	./"$(BACKEND_OUT)"

dev:
	@set -euo pipefail; \
	if [ -f .env ]; then set -a && . ./.env && set +a; fi; \
	( cd "$(WEB_DIR)" && npm install && npm run dev ) & \
	FRONT_PID="$$!"; \
	trap 'kill "$$FRONT_PID" 2>/dev/null || true' EXIT; \
	go run .; \
	wait "$$FRONT_PID"

clean:
	@rm -rf "$(BACKEND_OUT)" "$(WEB_ASSETS_DIR)"/*

docker:
        echo $PATH
        which docker
        ls -l /usr/bin/docker
        docker version
	docker build -t $(IMG_NAME):$(VERSION) -t $(IMG_NAME):latest -f Dockerfile .
