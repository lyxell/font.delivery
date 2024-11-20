build:
	go run main.go fonts_public.pb.go

test:
	go test -coverprofile=coverage.out

fmt:
  gofumpt -w .
