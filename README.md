# resource-aware-jds
Resource Aware Job Distribution

## Linter & Precommit Hook
1. Install golangci-lint
   https://github.com/golangci/golangci-lint
2. Install pre-commit (You may use brew)
   https://pre-commit.com/#install
3. Run
```bash
$ pre-commit install
```
4. Commit Code as normal or if you want to manually run lint you can use
```bash
$ make run
```

## GRPC Prerequisite
1. Install Protoc
```bash
$ brew install protobuf
```
2. Install protoc-gen-go
```bash
$ go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28
$ go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2
```
