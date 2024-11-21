compile:
	go build

build: compile
	rm -rf dist

	# Generate font files
	./fontdelivery --input-dir=fonts --output-dir=dist

	# Generate a master css file containing all font css files
	cat dist/*.css > dist/_.css

	# Copy static files to dist/
	#cp -r static/* dist/

	# Generate dist/style.css
	#tailwindcss -i input.css -o dist/style.css

typecheck:
	tsc -p jsconfig.json

test:
	go test -coverprofile=coverage.out

serve:
	miniserve --index index.html out/

fmt:
	gofumpt -w .
	prettier -w static/
