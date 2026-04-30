cover:
	go test ./... -coverprofile=coverage.out -covermode=count
	go tool cover -html=coverage.out -o coverage.html
	xdg-open coverage.html