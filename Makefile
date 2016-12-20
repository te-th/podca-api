# compile
.PHONY: compile
compile:
	@go build ./...

# run go metalinter
.PHONY: lint
lint:	compile
	@gometalinter --disable-all --enable=structcheck --enable=aligncheck --enable=deadcode --enable=gocyclo --enable=ineffassign --enable=dupl --enable=golint --enable=goimports --enable=varcheck --enable=interfacer --enable=goconst --enable=gosimple --enable=staticcheck --enable=unused --enable=misspell --enable=lll --line-length=120 --deadline=30s ./...

# run tests
.PHONY: test
test:	lint
	@go test $$(go list ./... | grep -v /vendor/)

# run build
.PHONY: build
build: test
	@goapp build ./app

# serve with goapp
.PHONY: serve
serve:
	@goapp serve ./app
