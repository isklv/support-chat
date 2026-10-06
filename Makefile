.PHONY: all build build-widget build-backend test dev clean docker

all: build

build: build-widget build-backend

build-widget:
	cd widget && pnpm install && pnpm build

build-backend:
	cd backend && go build -o bin/server ./cmd/server

test:
	cd backend && go test -v ./internal/...

dev-backend:
	cd backend && go run ./cmd/server

dev-widget:
	cd widget && pnpm dev

docker:
	docker compose up --build -d

clean:
	rm -rf backend/bin widget/dist
