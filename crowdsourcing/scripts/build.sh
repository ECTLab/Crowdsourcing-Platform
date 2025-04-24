#!/bin/sh

project_path=$(realpath $(dirname $(realpath $0))/..)

mkdir -p $project_path/release

go generate $project_path/...
go build -o $project_path/release $project_path/...
