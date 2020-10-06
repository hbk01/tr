main:
	@go build -o ./bin/tr
	@./bin/tr hello
install:
	go build -o ./bin/tr
	sudo cp ./bin/tr /usr/bin/tr
clean:
	rm -rf ./bin
