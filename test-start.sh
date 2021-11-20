#!/bin/sh

export TEST_MODE="TRUE"

cd utils
go test -v
cd ../
cd riotinterface
go test -v
cd ../