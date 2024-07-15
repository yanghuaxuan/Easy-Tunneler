#!/bin/sh

pushd backend
go run . &
popd

pushd frontend
npm run dev &
pushd