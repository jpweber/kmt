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

	"path/filepath"

	"github.com/spf13/viper"
)

var buildNumber string
var appVersion string

var paramList CLIParameters

func main() {

	// cli options

	versionPtr := flag.Bool("version", false, "Show version")
	filePath := flag.String("i", ".", "template file to input")
	// paramsFile := flag.String("f", "", "Parameter Values file rather than cli args. ")
	flag.Var(&paramList, "p", "<NAME>=<VALUE> Supplies a value for the named parameter")

	// Once all flags are declared, call `flag.Parse()`
	// to execute the command-line parsing.
	flag.Parse()

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
	err := viper.ReadInConfig()  // Find and read the config file
	if err != nil {              // Handle errors reading the config file
		panic(fmt.Errorf("Fatal error config file: %s", err))
	}

	// print the version
	if *versionPtr == true {
		fmt.Println(appVersion + " Build " + buildNumber)
		os.Exit(0)
	}

	// get the values from the values file in to a nice map
	var valuesFromFile = make(map[string]string)
	keys := viper.AllKeys()
	for _, key := range keys {
		valuesFromFile[key] = viper.Get(key).(string)
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
	fmt.Println(manifest)
}
