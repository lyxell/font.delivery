compile:
	go build

build-fonts: compile
	rm -rf dist
	# Generate font files
	./fontdelivery --input-dir=fonts --output-dir=dist
	# Generate a master css file containing all font css files
	cat dist/*.css > dist/_.css

build: build-fonts

test:
	go test -coverprofile=coverage.out

fmt:
	gofumpt -w .
