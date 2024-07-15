#!/bin/bash

set -e

# build directory
if [[ ! -e build ]]; then
    mkdir build
fi

# Build frontend
pushd frontend
npx vite build --outDir ../backend/public
popd

# Build backend
pushd backend
go build -o easy_tunneler
cp easy_tunneler ../build
cp -R ./public ../build/
popd