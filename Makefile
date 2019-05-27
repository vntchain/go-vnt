# This Makefile is meant to be used by people that do not usually work
# with Go source code. If you know what GOPATH is then you probably
# don't need to bother with make.

.PHONY: gvnt android ios gvnt-cross swarm all test clean
.PHONY: gvnt-linux gvnt-linux-386 gvnt-linux-amd64 gvnt-linux-mips64 gvnt-linux-mips64le
.PHONY: gvnt-linux-arm gvnt-linux-arm-5 gvnt-linux-arm-6 gvnt-linux-arm-7 gvnt-linux-arm64
.PHONY: gvnt-darwin gvnt-darwin-386 gvnt-darwin-amd64
.PHONY: gvnt-windows gvnt-windows-386 gvnt-windows-amd64

GOBIN = $(shell pwd)/build/bin
GO ?= latest

gvnt:
	build/env.sh go run build/ci.go install ./cmd/gvnt
	@echo "Done building."
	@echo "Run \"$(GOBIN)/gvnt\" to launch gvnt."

swarm:
	build/env.sh go run build/ci.go install ./cmd/swarm
	@echo "Done building."
	@echo "Run \"$(GOBIN)/swarm\" to launch swarm."

all:
	build/env.sh go run build/ci.go install

android:
	build/env.sh go run build/ci.go aar --local
	@echo "Done building."
	@echo "Import \"$(GOBIN)/gvnt.aar\" to use the library."

ios:
	build/env.sh go run build/ci.go xcode --local
	@echo "Done building."
	@echo "Import \"$(GOBIN)/Gvnt.framework\" to use the library."

test: all
	build/env.sh go run build/ci.go test

lint: ## Run linters.
	build/env.sh go run build/ci.go lint

clean:
	rm -fr build/_workspace/pkg/ $(GOBIN)/*

# The devtools target installs tools required for 'go generate'.
# You need to put $GOBIN (or $GOPATH/bin) in your PATH to use 'go generate'.

devtools:
	env GOBIN= go get -u golang.org/x/tools/cmd/stringer
	env GOBIN= go get -u github.com/kevinburke/go-bindata/go-bindata
	env GOBIN= go get -u github.com/fjl/gencodec
	env GOBIN= go get -u github.com/golang/protobuf/protoc-gen-go
	env GOBIN= go install ./cmd/abigen
	@type "npm" 2> /dev/null || echo 'Please install node.js and npm'
	@type "solc" 2> /dev/null || echo 'Please install solc'
	@type "protoc" 2> /dev/null || echo 'Please install protoc'

# Cross Compilation Targets (xgo)

gvnt-cross: gvnt-linux gvnt-darwin gvnt-windows gvnt-android gvnt-ios
	@echo "Full cross compilation done:"
	@ls -ld $(GOBIN)/gvnt-*

gvnt-linux: gvnt-linux-386 gvnt-linux-amd64 gvnt-linux-arm gvnt-linux-mips64 gvnt-linux-mips64le
	@echo "Linux cross compilation done:"
	@ls -ld $(GOBIN)/gvnt-linux-*

gvnt-linux-386:
	build/env.sh go run build/ci.go xgo -- --go=$(GO) --targets=linux/386 -v ./cmd/gvnt
	@echo "Linux 386 cross compilation done:"
	@ls -ld $(GOBIN)/gvnt-linux-* | grep 386

gvnt-linux-amd64:
	build/env.sh go run build/ci.go xgo -- --go=$(GO) --targets=linux/amd64 -v ./cmd/gvnt
	@echo "Linux amd64 cross compilation done:"
	@ls -ld $(GOBIN)/gvnt-linux-* | grep amd64

gvnt-linux-arm: gvnt-linux-arm-5 gvnt-linux-arm-6 gvnt-linux-arm-7 gvnt-linux-arm64
	@echo "Linux ARM cross compilation done:"
	@ls -ld $(GOBIN)/gvnt-linux-* | grep arm

gvnt-linux-arm-5:
	build/env.sh go run build/ci.go xgo -- --go=$(GO) --targets=linux/arm-5 -v ./cmd/gvnt
	@echo "Linux ARMv5 cross compilation done:"
	@ls -ld $(GOBIN)/gvnt-linux-* | grep arm-5

gvnt-linux-arm-6:
	build/env.sh go run build/ci.go xgo -- --go=$(GO) --targets=linux/arm-6 -v ./cmd/gvnt
	@echo "Linux ARMv6 cross compilation done:"
	@ls -ld $(GOBIN)/gvnt-linux-* | grep arm-6

gvnt-linux-arm-7:
	build/env.sh go run build/ci.go xgo -- --go=$(GO) --targets=linux/arm-7 -v ./cmd/gvnt
	@echo "Linux ARMv7 cross compilation done:"
	@ls -ld $(GOBIN)/gvnt-linux-* | grep arm-7

gvnt-linux-arm64:
	build/env.sh go run build/ci.go xgo -- --go=$(GO) --targets=linux/arm64 -v ./cmd/gvnt
	@echo "Linux ARM64 cross compilation done:"
	@ls -ld $(GOBIN)/gvnt-linux-* | grep arm64

gvnt-linux-mips:
	build/env.sh go run build/ci.go xgo -- --go=$(GO) --targets=linux/mips --ldflags '-extldflags "-static"' -v ./cmd/gvnt
	@echo "Linux MIPS cross compilation done:"
	@ls -ld $(GOBIN)/gvnt-linux-* | grep mips

gvnt-linux-mipsle:
	build/env.sh go run build/ci.go xgo -- --go=$(GO) --targets=linux/mipsle --ldflags '-extldflags "-static"' -v ./cmd/gvnt
	@echo "Linux MIPSle cross compilation done:"
	@ls -ld $(GOBIN)/gvnt-linux-* | grep mipsle

gvnt-linux-mips64:
	build/env.sh go run build/ci.go xgo -- --go=$(GO) --targets=linux/mips64 --ldflags '-extldflags "-static"' -v ./cmd/gvnt
	@echo "Linux MIPS64 cross compilation done:"
	@ls -ld $(GOBIN)/gvnt-linux-* | grep mips64

gvnt-linux-mips64le:
	build/env.sh go run build/ci.go xgo -- --go=$(GO) --targets=linux/mips64le --ldflags '-extldflags "-static"' -v ./cmd/gvnt
	@echo "Linux MIPS64le cross compilation done:"
	@ls -ld $(GOBIN)/gvnt-linux-* | grep mips64le

gvnt-darwin: gvnt-darwin-386 gvnt-darwin-amd64
	@echo "Darwin cross compilation done:"
	@ls -ld $(GOBIN)/gvnt-darwin-*

gvnt-darwin-386:
	build/env.sh go run build/ci.go xgo -- --go=$(GO) --targets=darwin/386 -v ./cmd/gvnt
	@echo "Darwin 386 cross compilation done:"
	@ls -ld $(GOBIN)/gvnt-darwin-* | grep 386

gvnt-darwin-amd64:
	build/env.sh go run build/ci.go xgo -- --go=$(GO) --targets=darwin/amd64 -v ./cmd/gvnt
	@echo "Darwin amd64 cross compilation done:"
	@ls -ld $(GOBIN)/gvnt-darwin-* | grep amd64

gvnt-windows: gvnt-windows-386 gvnt-windows-amd64
	@echo "Windows cross compilation done:"
	@ls -ld $(GOBIN)/gvnt-windows-*

gvnt-windows-386:
	build/env.sh go run build/ci.go xgo -- --go=$(GO) --targets=windows/386 -v ./cmd/gvnt
	@echo "Windows 386 cross compilation done:"
	@ls -ld $(GOBIN)/gvnt-windows-* | grep 386

gvnt-windows-amd64:
	build/env.sh go run build/ci.go xgo -- --go=$(GO) --targets=windows/amd64 -v ./cmd/gvnt
	@echo "Windows amd64 cross compilation done:"
	@ls -ld $(GOBIN)/gvnt-windows-* | grep amd64
