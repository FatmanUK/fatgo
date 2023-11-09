#!/bin/bash
set -euo pipefail


cd serialisers
go test -v
cd ..


for D in *; do
	[ ! -d ${D} ] && continue
	pushd ${D} >/dev/null
	go test -v && true
	popd >/dev/null
done

