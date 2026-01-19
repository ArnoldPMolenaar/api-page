.PHONY: clean critic security lint

APP_NAME = api-page
BUILD_DIR = $(PWD)/build

clean:
	rm -rf ./build

critic:
	gocritic check -enableAll ./...

lint:
	golangci-lint run ./...
