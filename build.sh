#!/bin/bash

set -e

# build directory
if [[ ! -e build ]]; then
    mkdir build
fi

# Build frontend
pushd frontend
npm i
npx vite build --outDir ../backend/public
popd

# Build backend
pushd backend
go build -o easy_tunneler
mv easy_tunneler ../build
if [[ -e ../build/public ]]; then
    rm -Rf ../build/public
fi
mv ./public ../build/
popd

echo 'EASY_TUNNELER_PROD=1 GIN_MODE=release ./easy_tunneler' > build/run.sh; chmod +x build/run.sh
