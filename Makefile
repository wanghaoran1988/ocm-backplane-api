build: clean
	GO111MODULE=on go build -mod=vendor -o ocm-backplane-api

clean:
	rm -f ocm-backplane-api

.PHONY: lint
lint:
	golangci-lint run
