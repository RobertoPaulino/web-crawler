.PHONY: build test run clean web

build:
	go build -o bin/crawler ./cmd/crawler
	go build -o bin/web ./cmd/web

test:
	go test -v ./...

run: build
	./bin/crawler $(ARGS)

web: build
	./bin/web

clean:
	rm -rf bin/

# Example usage:
# make run ARGS="https://example.com 10 100"
# make web (to start the web interface on port 8080) 