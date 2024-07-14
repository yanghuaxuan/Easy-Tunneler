#!/bin/sh

# build directory
if [[ ! -e build ]]; then
    mkdir build
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
cp -R easy_tunneler ../build
cp -R ./public ../build/
popd
