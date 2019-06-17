#!/bin/sh

cd /go/src/github.com/nalej/system-model
go test -v -coverprofile=coverage.out -covermode=atomic ./... 2>&1
echo "---####--- COVERAGE FILE CONTENTS ---####---"
cat coverage.out
