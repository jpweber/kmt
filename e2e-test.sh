#!/bin/bash
appname="./kmt"
go build

## Test only reading values from file
# result=$(./$appname -i testfiles/test.tmpl)
# testfile=$(cat testfiles/test-result.yaml)

# if [ "$result" != "$testfile" ] 
# then
#     echo "Files do not match"
#     echo "$(diff <(echo "$result") <(echo "$testfile"))"
# else
#     echo "File Only Test Good!"
# fi

## test reading values from file and from cli.
## cli values should override file
# result=$(./$appname -i testfiles/test.tmpl -p version=2.2.2 -p name=testcontainer2)
# testfile=$(cat testfiles/test-result2.yaml)
# if [ "$result" != "$testfile" ] 
# then
#     echo "Files do not match"
#     echo "$(diff <(echo "$result") <(echo "$testfile"))"
# else
#     echo "File and CLI params Test Good!"
# fi

## test that we properly skip making artifacts if values not present.
## cli values should override file
result=$(./$appname -i testfiles/test-noarts.tmpl -p namespace=prod -v)
testfile=$(cat testfiles/test-result-noarts.yaml)
if [ "$result" != "$testfile" ] 
then
    echo "Files do not match"
    echo "$(diff <(echo "$result") <(echo "$testfile"))"
else
    echo "File and CLI params Test Good!"
fi
