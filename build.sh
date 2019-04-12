#!/bin/bash

tool_name=("findcpp" "codecount")
dest_dir="tools"
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
    go build $tool.go
    if [ -f $tool ]
    then
        mv $tool $tool_dir
    else
        echo "build failed: $tool"
    fi
done