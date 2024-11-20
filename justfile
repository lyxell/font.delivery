build:
	rm -rf out
	go run main.go fonts_public.pb.go
	cat out/*.css > out/_.css
	cp -r static/* out/
	# Generate out/style.css
	tailwindcss -i input.css -o out/style.css

typecheck:
	tsc -p jsconfig.json

test:
	go test -coverprofile=coverage.out

serve:
	miniserve --index index.html out/

fmt:
	gofumpt -w .
	prettier -w static/
