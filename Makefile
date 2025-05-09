build:
	go build -v ./...

install: build
	go install -v ./...

fmt:
	gofmt -s -w -e .

testacc:
	TF_ACC=1 go test -v -cover -timeout 120m ./...

generate:
	cd tools; go generate ./...