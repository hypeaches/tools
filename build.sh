#!/bin/bash

tool_name=("codecount" "dig" "findcpp")
dest_dir="tools"

# 构建脚本
script_dir=$(cd $(dirname ${BASH_SOURCE[0]}); pwd)

tool_dir=$script_dir/$dest_dir
if [ -d $tool_dir ]
then
    rm -rf $tool_dir
fi
mkdir $tool_dir

for tool in ${tool_name[@]}
do
    cd $script_dir/$tool
    echo "build: $tool"
    go build $tool.go
    if [ -f $tool ]
    then
        mv $tool $tool_dir
    else
        echo "build failed: $tool"
    fi
done
