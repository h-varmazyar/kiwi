#!/bin/sh
echo "running ${1} ${2}"
cd ./applications/"${1}"/cmd && go run . "${FLAGS}" "${2}"