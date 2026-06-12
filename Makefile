GOBIN := $(shell go env GOPATH)/bin

.PHONY: fmt fmt-check vet lint check

fmt:
	cd backend && gofmt -w . && $(GOBIN)/goimports -w -local backend .
	cd frontend && npx prettier --write "src/**/*.{ts,html,scss}"

fmt-check:
	cd backend && test -z "$$(gofmt -l .)" && test -z "$$($(GOBIN)/goimports -l -local backend .)"

vet:
	cd backend && go vet ./...

lint:
	cd backend && $(GOBIN)/golangci-lint run
	cd frontend && npm run lint

check: fmt-check vet lint
