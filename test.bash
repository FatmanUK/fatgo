#!/bin/bash
set -euo pipefail

for D in *; do
	[ ! -d ${D} ] && continue
	pushd ${D} >/dev/null
	go test -v && true
	popd >/dev/null
done
