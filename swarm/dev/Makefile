.PHONY: build cluster test

default: build

build:
	go build -o bin/swarm    github.com/vntchain/go-vnt/cmd/swarm
	go build -o bin/gvnt     github.com/vntchain/go-vnt/cmd/gvnt
	go build -o bin/bootnode github.com/vntchain/go-vnt/cmd/bootnode

cluster: build
	scripts/boot-cluster.sh

test:
	go test -v github.com/vntchain/go-vnt/swarm/...
