# run go metalinter
.PHONY: lint
lint:
	@gometalinter ./...

# compile
.PHONY: compile
compile:
	@go build ./...

# run tests
.PHONY: test
test:
	@go test $$(go list ./... | grep -v /vendor/ | grep -v /cmd/)

# run build
.PHONY: build
build: compile
	@goapp build ./app

# serve with goapp
.PHONY: serve
serve:
	@goapp serve ./app
