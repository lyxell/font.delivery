build:
	rm -rf out
	go run main.go fonts_public.pb.go
	cat out/*.css > out/_.css
	cp -r static/* out/

test:
	go test -coverprofile=coverage.out

fmt:
	gofumpt -w .
