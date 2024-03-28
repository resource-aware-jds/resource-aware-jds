#!/bin/bash

go build -o out/cp ./cmd/controlplane/main.go

go build -o out/worker ./cmd/worker/main.go
