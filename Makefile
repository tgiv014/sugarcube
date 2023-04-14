.PHONY: dev
## dev: starts app in dev mode
dev:
	go run cmd/dev/main.go

.PHONY: build
## build: builds server with embedded frontend
build:
	cd web && npm i && npm run build
	go build -o sugarcube cmd/sugarcube/main.go
