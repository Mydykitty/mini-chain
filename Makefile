
build:
	CGO_ENABLED=0 go build \
		main.go

vendor:tidy
	go mod vendor
tidy:
	go mod tidy
