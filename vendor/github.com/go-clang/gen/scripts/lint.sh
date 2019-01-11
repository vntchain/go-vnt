 #!/bin/bash

echo "gofmt:"
gofmt -l . 2>&1 | grep --invert-match -P "(_gen.go|testdata/)" || true

echo "goerrcheck:"
errcheck ./... 2>&1 | grep --invert-match -P "(_gen.go|testdata/)" || true

echo "govet:"
go tool vet -all=true -v=true . 2>&1 | grep --invert-match -P "(_gen.go|testdata/|Checking file|\%p of wrong type|can't check non-constant format)" || true

echo "golint:"
golint ./... 2>&1 | grep --invert-match -P "(_gen.go|testdata/|_string.go:)" || true
