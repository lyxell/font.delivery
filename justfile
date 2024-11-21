compile:
	go build

build-frontend:
	cd frontend && just build

build-fonts: compile
	rm -rf dist
	# Generate font files
	./fontdelivery --input-dir=fonts --output-dir=dist
	# Generate a master css file containing all font css files
	cat dist/*.css > dist/_.css

build: build-fonts

serve:
	cd frontend && just serve

test:
	go test -coverprofile=coverage.out

fmt:
	gofumpt -w .
