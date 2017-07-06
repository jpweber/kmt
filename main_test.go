package main

import (
	"io/ioutil"
	"testing"
)

func TestparseManifestTmpl(t *testing.T) {

	paramList := []string{"version=1.0.0", "name=testcontainer"}
	// init Parameters map
	// parameters := ManifestValues{}
	parameters := ManifestValues{
		Values: initParams(paramList),
	}

	filePath := "./testfiles/test.yaml"

	// read in the tmplate file
	tmplBytes, _ := ioutil.ReadFile(filePath)
	manifestTmpl := string(tmplBytes)
	manifest := parseManifestTmpl(parameters, manifestTmpl)
	testManifest := `
	apiVersion: v1
	kind: Pod
	metadata:
	name: testcontainer
	namespace: default
	spec:
	containers:
	- image: busybox:1.0.0
		command:
		- sleep
		- "3600"
		imagePullPolicy: IfNotPresent
		name: busybox
	restartPolicy: Always`

	if manifest != testManifest {
		t.Error("Manifest files do not match got ", testManifest)
	}
}
