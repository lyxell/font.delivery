compile-builder:
	go build ./cmd/builder

compile-fontdl:
	oapi-codegen -generate client,types -o internal/api/api.go api.yml
	go build ./cmd/fontdl

# Build web interface
build-web:
	cd web && just build
	cp -r web/dist/* dist/
	echo "404 Not Found" > dist/404.html

# Generate WOFF2 and CSS files
build-fonts: compile-builder
	./builder --input-dir=fonts/ --output-dir=dist/
	# Generate a master css file containing all font css files
	cat dist/*.css > dist/api/v1/download/_.css

# Build API docs
build-api-docs:
	redocly build-docs --output=dist/reference/index.html api.yml

clean-dist:
	rm -rf dist/*

build: clean-dist build-fonts build-api-docs build-web

serve:
	cd web && just serve

serve-production:
	miniserve --index index.html dist/

lint-api:
	redocly lint api.yml

test:
	go test -coverprofile=coverage.out ./...

fmt:
	gofumpt -w .
