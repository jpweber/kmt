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
	"log"
	"os"

	"path/filepath"

	"github.com/spf13/viper"
)

var buildNumber string
var appVersion = "1.2.2"
var debug = false

var paramList CLIParameters

func logger(logMsg string) {
	if debug {
		log.Println(logMsg)
	}
}

func main() {

	// cli options

	versionPtr := flag.Bool("version", false, "Show version")
	filePath := flag.String("i", ".", "template file to input")
	// paramsFile := flag.String("f", "", "Parameter Values file rather than cli args. ")
	flag.Var(&paramList, "p", "<NAME>=<VALUE> Supplies a value for the named parameter")
	verbose := flag.Bool("v", false, "Print Parsed Templated to STDOUT")
	xtraVerbose := flag.Bool("vv", false, "Print Parsed Templated to STDOUT Plus log messages. This is not good for piping to kubectl ")

	// Once all flags are declared, call `flag.Parse()`
	// to execute the command-line parsing.
	flag.Parse()

	// print the version
	if *versionPtr == true {
		fmt.Println(appVersion)
		os.Exit(0)
	}

	// check for extra verbose output
	if *xtraVerbose == true {
		debug = true
	}

	// get absolute path to manifest file
	manPath, _ := filepath.Abs(*filePath)
	// get the filename itself
	_, manFile := filepath.Split(manPath)
	// get just the path to the file, excluding the file itself
	manPath = filepath.Dir(manPath)

	extension := filepath.Ext(manFile)
	noExtFileName := manFile[0 : len(manFile)-len(extension)]

	viper.SetConfigName(noExtFileName + "-values") // name of config file (without extension)
	// viper.AddConfigPath("/etc/appname/")  // path to look for the config file in
	// viper.AddConfigPath("$HOME/.appname") // call multiple times to add many search paths
	viper.AddConfigPath(manPath) // optionally look for config in the working directory

	err := viper.ReadInConfig() // Find and read the config file
	if err != nil {             // Handle errors reading the config file
		// check if the config file simply doesn't exist
		// it is not a requirement so move on if that is the error
		_, err := os.Stat(viper.ConfigFileUsed())
		if err != nil {
			logger("Values file does not exist. Not required. Moving on.")
		} else {
			// any other errors reading in the config should cause us to stop
			panic(fmt.Errorf("Fatal error config file: %s", err))
		}

	}

	// get the values from the values file in to a nice map
	var valuesFromFile = make(map[string]interface{})
	keys := viper.AllKeys()
	for _, key := range keys {
		// need to build up list of other value types to infer correctly

		valuesFromFile[key] = viper.Get(key)

	}

	// get the params passed as CLI args in to a  nice map
	valuesFromCLI := initParams(paramList)

	// Add any CLI values to the
	parameters := ManifestValues{
		Values: mergeParams(valuesFromFile, valuesFromCLI),
	}

	// read in the tmplate file
	tmplBytes, _ := ioutil.ReadFile(*filePath)
	manifestTmpl := string(tmplBytes)
	manifest := parseManifestTmpl(parameters, manifestTmpl)

	// print file output to screen if verbose
	if *verbose {
		fmt.Println(manifest)
	}

	if parameters.Values["name"] != nil && parameters.Values["version"] != nil {
		logger("writing artifacts")
		writeArtifacts(manPath, noExtFileName, manifest, parameters)
	}

}
