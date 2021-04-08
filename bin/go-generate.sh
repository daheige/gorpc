#!/usr/bin/env bash
root_dir=$(cd "$(dirname "$0")"; cd ..; pwd)

protoExec=$(which "protoc")
if [ -z $protoExec ]; then
    echo 'Please install protoc!'
    echo "Please look readme.md to install proto3"
    echo "if you use centos7,please look https://github.com/daheige/go-proj/blob/master/docs/centos7-protoc-install.md"
    exit 0
fi

# 安装grpc_tools
goInjectExec=$(which "protoc-go-inject-tag")
if [ -z goInjectExec ]; then
  sh $root_dir/bin/grpc_tools.sh
fi

proto_dir=$root_dir/api/protos
pb_dir=$root_dir/api/clients/go/pb

mkdir -p $pb_dir
mkdir -p $proto_dir

#delete old pb code.
rm -rf $pb_dir/*.go

echo "\n\033[0;32mGenerating codes...\033[39;49;0m\n"

echo "generating golang stubs..."

# 生成grpc代码和http gateway代码
$protoExec -I $proto_dir --go_out=plugins=grpc:$pb_dir --grpc-gateway_out=logtostderr=true:$pb_dir $proto_dir/*.proto

#inject tag
echo "\n\033[0;32minject tag...\033[39;49;0m\n"
sh $root_dir/bin/protoc-inject-tag.sh

# request validator code
sh $root_dir/bin/validator-generate.sh

echo "generating golang code success"
echo "\n\033[0;32mGenerate codes successfully!\033[39;49;0m\n"

exit 0
