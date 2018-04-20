#!/bin/sh

dir=$(dirname "$0")
srcdir=$dir/..

GOOS=linux go build -o "$srcdir/handler" "$srcdir"
