/*
* @Author: Jim Weber
* @Date:   2017-04-23 22:57:47
* @Last Modified by:   Jim Weber
* @Last Modified time: 2017-04-23 23:10:31
 */

package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
)

var buildNumber string
var appVersion string
var paramList CLIParameters

func main() {

	// cli options

	versionPtr := flag.Bool("version", false, "Show version")
	filePath := flag.String("i", ".", "template file to input")
	flag.Var(&paramList, "p", "<NAME>=<VALUE> Supplies a value for the named parameter")

	// Once all flags are declared, call `flag.Parse()`
	// to execute the command-line parsing.
	flag.Parse()

	if *versionPtr == true {
		fmt.Println(appVersion + " Build " + buildNumber)
		os.Exit(0)
	}

	// init Parameters map
	// parameters := ManifestValues{}
	parameters := ManifestValues{
		Values: initParams(paramList),
	}

	fmt.Println(parameters.Values["name"])

	// read in the tmplate file
	tmplBytes, _ := ioutil.ReadFile(*filePath)
	manifestTmpl := string(tmplBytes)

	manifest := parseManifestTmpl(parameters, manifestTmpl)
	fmt.Println(manifest)
}
