# This action is run by CI when releasing fontdl
build-fontdl NAME:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o bin/{{NAME}}-linux-amd64 ./cmd/fontdl
	CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -o bin/{{NAME}}-macos-amd64 ./cmd/fontdl
	CGO_ENABLED=0 GOOS=darwin GOARCH=arm64 go build -o bin/{{NAME}}-macos-arm64 ./cmd/fontdl

package-fontdl:
	just build-fontdl "$(git describe --exact-match --tags)"

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
