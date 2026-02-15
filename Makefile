BINARY := app
PORT   := 8080
IMAGE  := app:latest
MAIN   := .

LDFLAGS := -s -w

.PHONY: tidy fmt vet test build run clean docker-build docker-run docker-up docker-down

tidy:
	go mod tidy

fmt:
	go fmt ./...

vet:
	go vet ./...

test:
	go test ./...

build:
	mkdir -p bin
	CGO_ENABLED=0 go build -ldflags="$(LDFLAGS)" -o bin/$(BINARY) $(MAIN)

run:
	go run $(MAIN)

clean:
	rm -rf bin

docker-build:
	docker build -t $(IMAGE) .

docker-run:
	docker run --rm -p $(PORT):$(PORT) --name $(BINARY) $(IMAGE)

docker-up:
	docker run -d -p $(PORT):$(PORT) --name $(BINARY) $(IMAGE)

docker-down:
	docker rm -f $(BINARY) || true
