#!/bin/bash

set -e

# build directory
if [[ ! -e build ]]; then
    mkdir build
fi

if [[ ! -e .env ]]; then
    echo 'Cannot find .env! Exiting.'
    exit 1
fi

# Copy .env to frontend directory for Vite to build with
awk '{sub(/^[A-Z_]*=/, "VITE_&"); print}' .env > ./frontend/.env

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
