#!/bin/bash

PROG=mocker

go build -o $PROG ./...

rm -rf out
./$PROG gen android

pushd out > /dev/null
./gradlew assemble
popd > /dev/null

rm -f $PROG

