#!/bin/bash

appname="kmt"
go build

## Test only reading values from file
result=$(./$appname -i testfiles/test.yaml)
testfile=$(cat testfiles/test-result.yaml)

if [ "$result" != "$testfile" ] 
then
    echo "Files do not match"
    echo "$(diff <(echo "$result") <(echo "$testfile"))"
else
    echo "File Only Test Good!"
fi

## test reading values from file and from cli.
## cli values should override file
result=$(./$appname -i testfiles/test.yaml -p version=2.2.2 -p name=testcontainer2)
testfile=$(cat testfiles/test-result2.yaml)
if [ "$result" != "$testfile" ] 
then
    echo "Files do not match"
    echo "$(diff <(echo "$result") <(echo "$testfile"))"
else
    echo "File and CLI params Test Good!"
fi

