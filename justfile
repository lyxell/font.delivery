test:
	go test -coverprofile=coverage.out

fmt:
  gofumpt -w .
