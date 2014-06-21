bin/henchman: fetch-deps
	go build -x -o $@ ./

clean:
	rm -rf bin/

all: bin/henchman

fetch-deps:
	go get -d -v ./...

.PHONEY: clean all fetch-deps
