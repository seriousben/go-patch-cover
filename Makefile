
testdata/test-project/coverage.out : $(wildcard testdata/test-project/*.go)
	rm $@ \
	cd testdata/test-project \
	go test -coverprofile=coverage.out -covermode=count ./...
