rebuild:
	make clean;
	make build;
	./validator

clean:
	rm -rf validator

build:
	gofmt -w .
	go build -o validator