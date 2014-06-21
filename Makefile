test: fetch-deps
	go test ./...

clean:
	rm -rf bin/

bin/henchman: fetch-deps
	go build -x -o $@ ./

all: bin/henchman

fetch-deps:
	go get -d -v ./...

.PHONEY: clean all test fetch-deps
