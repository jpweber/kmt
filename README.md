# KMT

## Description
kmt is short for Kubernetes Manifest Template. It is a tool for creating templated manifest files
and then populating certain values prorgramitcally.

It supports two main ways of dynamically populating data.
- 1) pass arguments via the cli `-p foo=bar` would populate the place holder of foo wit the values "bar"

- 2) Another way is a values file. This file is name just like your template file but with `-values.yaml` at the end. kmt will automatically look for these files to populate values. If it does not exist it will then use cli args. 

You can use both a values file and cli args together. the cli parameters will override any thing in the values file.

kmt will also try to save manifest artifacts if a name and version are used as parameters either via cli or a values file. The path it will save these fully composed manifest files is as follows. Artifacts will be placed in the root dir of the manifest file itself.  If nothing that maps to a namespace is provided `default` will be used as the namespace
`./artifacts/<namespace>/<app name>/filename`  
a dir container latest versions of the file is also use
`./artifacts/<namespace>/current/filename`


## Usage
```
Usage of ./kmt:
  -f string
    	Parameter Values file rather than cli args.
  -i string
    	template file to input
  -p value
    	<NAME>=<VALUE> Supplies a value for the named parameter
  -v	Print Parsed Templated to STDOUT
  -version
    	Show version
  -vv
    	Print Parsed Templated to STDOUT Plus log messages. This is not good for piping to kubectl
```

the verbose option will output the composed manifest file to `STDOUT` whigh is usful for piping to other commands, such as `kubectl` for example.