
build:
	clear
	go build -o ./bin/dooms-store.bin main.go


run: build
	./bin/dooms-store.bin

test: build
	go test ./... -v --race