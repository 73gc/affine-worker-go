#!/bin/bash
CURDIR=$(cd $(dirname $0); pwd)
BinaryName=affine_worker_go
echo "$CURDIR/bin/${BinaryName}"
exec $CURDIR/bin/${BinaryName}