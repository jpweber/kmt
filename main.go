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
	"time"

	"path/filepath"

	"github.com/spf13/viper"
)

var buildNumber string
var appVersion = "1.3.2"
var debug = false

var paramList CLIParameters

func logger(logMsg string) {
	if debug {
		log.Println(logMsg)
	}
}

func stripExt(filename string) string {
	extension := filepath.Ext(filename)
	noExtFileName := filename[0 : len(filename)-len(extension)]

	return noExtFileName
}

func main() {

	// cli options

	versionPtr := flag.Bool("version", false, "Show version")
	filePath := flag.String("i", "", "template file to input")
	paramsFile := flag.String("f", "", "Parameter Values file rather than cli args. ")
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

	// show usage if one of the required parameters is not provided.
	if *filePath == "" && *paramsFile == "" {
		flag.Usage()
		os.Exit(1)
	}

	// get absolute path to manifest file
	manPath, _ := filepath.Abs(*filePath)
	// get the filename itself
	_, manFile := filepath.Split(manPath)
	// get just the path to the file, excluding the file itself
	manPath = filepath.Dir(manPath)
	extension := filepath.Ext(manFile)
	noExtFileName := manFile[0 : len(manFile)-len(extension)]

	if *paramsFile != "" {
		// someone specified a values file. with the -f flag
		// Use this one instead of deriving one based on convention
		viper.SetConfigName(stripExt(*paramsFile))
		// reset manpath to the path to the paramsFile
		// for use when we read in the config file
		manPath, _ = filepath.Abs(*paramsFile)
		manPath = filepath.Dir(manPath)
	} else {
		// someone provided a normal template.
		// set the config file based on the convesion of templatename-values
		viper.SetConfigName(noExtFileName + "-values")
	}
	// look for config in the working directory
	viper.AddConfigPath(manPath)

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

	// Add in a timestamp to the values for anyone templates that
	// are expecting one.
	t := time.Now().UTC()
	parameters.Values["date"] = fmt.Sprintf(t.Format("20060102_15:04:05"))

	// If an template file was provided via CLI args read that in.
	var manifestTmpl string
	if *filePath != "" {
		tmplBytes, _ := ioutil.ReadFile(*filePath)
		manifestTmpl = string(tmplBytes)
	} else {
		logger("generating base template from provided template name in values file")
		manifestTmpl = "{{ template  \"" + parameters.Values["template"].(string) + "\" . }}"
	}

	// populate the template with values
	manifest := parseManifestTmpl(parameters, manifestTmpl)

	// print file output to screen if verbose
	if *verbose || *xtraVerbose {
		fmt.Println(manifest)
	}

	// save the finished file to disk in an artifacts dir.
	if parameters.Values["name"] != nil && parameters.Values["version"] != nil {
		logger("writing artifacts")
		writeArtifacts(manPath, noExtFileName, manifest, parameters)
	}

}
