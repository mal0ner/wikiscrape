VERSION:=0.1.0

OS = darwin freebsd linux openbsd
ARCHS = 386 arm amd64 arm64

all: build release

build: deps
	go build

release:
	@for arch in $(ARCHS);\
	do \
		for os in $(OS);\
		do \
			echo "Building $$os-$$arch"; \
			mkdir -p build/wikiscrape-$$os-$$arch/; \
			GOOS=$$os GOARCH=$$arch go build -o build/wikiscrape-$$os-$$arch/wikiscrape; \
			## tar cz -C build -f build/wikiscrape-$$os-$$arch.tar.gz wikiscrape-$$os-$$arch; \
		done \
	done \

test: deps
	go test ./...

deps:
	go get -d -v -t ./...

clean:
	rm -rf build
	rm -f wikiscrape
