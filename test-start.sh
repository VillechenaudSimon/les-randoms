#!/bin/sh

. ./setup-secret-keys.sh

export TEST_MODE="TRUE"
export DEBUG_MODE="FALSE"

cd utils
go test -v
cd ../
cd riotinterface
go test -v
cd ../