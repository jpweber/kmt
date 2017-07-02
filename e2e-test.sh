#!/bin/bash

appname="kmantmpl"
go build
result=$(./$appname -i testfiles/test.yaml)
testfile=$(cat testfiles/test-result.yaml)

if [ "$result" != "$testfile" ] 
then
    echo "Files do not match"
    echo "$(diff <(echo "$result") <(echo "$testfile"))"
else
    echo "Test Good!"
fi
