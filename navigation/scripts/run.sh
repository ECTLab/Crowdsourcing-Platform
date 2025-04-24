#!/bin/sh

project_path=$(realpath $(dirname $(realpath $0))/..)

go generate $project_path/...
go run $project_path/cmd/navigation/main.go
