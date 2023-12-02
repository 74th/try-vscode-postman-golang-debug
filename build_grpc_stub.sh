#!/bin/bash
set -x
mkdir -p ./server/grpc/pb
protoc --plugin=/home/nnyn/go/bin/protoc-gen-go --go_out=./server/grpc/pb --go_opt=paths=source_relative \
    --plugin=/home/nnyn/go/bin/protoc-gen-go-grpc --go-grpc_out=./server/grpc/pb --go-grpc_opt=paths=source_relative \
    gRPC.proto