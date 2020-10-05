main:
	@go build -o ./bin/tr
	@./bin/tr hello
clean:
	rm -rf ./bin
