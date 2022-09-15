test:
	go test -count=1 ./...

lint:
	golangci-lint run

kappctl-install:
	go install cmd/kappctl/kappctl.go
