test: fetch-deps
	gom test ./...

clean:
	rm -rf bin/
	rm -rf _vendor

bin/henchman: fetch-deps
	gom build -x -o $@ ./

all: bin/henchman

fetch-deps:
	go get github.com/mattn/gom
	gom install

.PHONEY: clean all test fetch-deps
