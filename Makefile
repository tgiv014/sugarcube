.PHONY: dev
## dev: starts app in dev mode
dev:
	ENVIRONMENT=dev go run cmd/sugarcube/main.go

.PHONY: build
## build: builds server with embedded frontend
build:
	cd web && npm run build
	go build cmd/sugarcube/main.go
