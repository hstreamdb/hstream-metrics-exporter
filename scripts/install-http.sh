#!/bin/sh
set -e

DEV_DIST="./_dev_dist"
mkdir -p $DEV_DIST

cd $DEV_DIST

git clone --depth=1 https://github.com/hstreamdb/http-services.git

cd http-services

make
