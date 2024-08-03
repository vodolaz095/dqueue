start:
	go run example/example.go

race:
# not very reliable, since it gives different results
	CGO_ENABLED=1 go run --race example/example.go

# installing golint code quality tools and checking, if it can be started
# go install golang.org/x/lint/golint@latest
lint:
	gofmt -w=true -s=true -l=true ./
	golint ./...
	go vet ./...

check: lint

test:
	go test -v ./...

cover:
	go test -v --cover ./...
