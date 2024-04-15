
build:
	clear
	go build -o ./bin/dooms-store.bin main.go


run: build
	./bin/dooms-store.bin

test: build
	go test ./... -v --race


run-container:
	clear
	podman ps
	@echo "removing the old docker images"
	podman images | grep dooms | awk '{print $$3}' | xargs podman rmi -f
	podman volume prune
	@echo "building latest docker image"
	podman build -t dooms-store -f Dockerfile .
	@echo "running the docker container"
	podman run -ti -p 3000:3000 dooms-store