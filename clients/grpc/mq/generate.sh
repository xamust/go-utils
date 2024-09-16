#!/bin/sh -ex

DIR_GEN="./generate"

rm -rf ${DIR_GEN}
mkdir ${DIR_GEN}
protoc -I ./ \
  --go_out ${DIR_GEN} --go_opt paths=source_relative \
  --go-grpc_out ${DIR_GEN} --go-grpc_opt paths=source_relative \
  --grpc-gateway_out ${DIR_GEN} --grpc-gateway_opt paths=source_relative \
  ./proto/*.proto \
  || rm -rf ${DIR_GEN}