/*
* @Author: Jim Weber
* @Date:   2017-04-23 23:28:02
* @Last Modified by:   Jim Weber
* @Last Modified time: 2017-04-23 23:28:06
 */

package main

import (
	"fmt"
	"strings"
)

// CLIParameters Define a type named "Parameters" as a slice of strings
type CLIParameters []string

// Now, for our new type, implement the two methods of
// the flag.Value interface...
// The first method is String() string
func (i *CLIParameters) String() string {
	return fmt.Sprintf("%d", *i)
}

// Set the extra var value Set(value string) error
func (i *CLIParameters) Set(value string) error {
	tmp := value
	*i = append(*i, tmp)
	return nil
}

func initParams(paramList []string) map[string]interface{} {
	parameters := make(map[string]interface{})
	for i := 0; i < len(paramList); i++ {
		varParts := strings.Split(paramList[i], "=")
		if len(varParts) > 1 {
			parameters[varParts[0]] = strings.Join(varParts[1:], ":")
		} else {
			fmt.Println("Error with", paramList[i], "Key value pair. Format should be key=value")
		}
	}
	return parameters
}

func mergeParams(fromFile, fromCLI map[string]interface{}) map[string]interface{} {
	// init the final parameter values list
	var finalParams = make(map[string]interface{})

	//Add the values from file to the final params
	if len(fromFile) > 0 {
		for k, v := range fromFile {
			finalParams[k] = v
		}
	}

	// Add the values from the CLI overriding any values from file
	if len(fromCLI) > 0 {
		for k, v := range fromCLI {
			finalParams[k] = v
		}
	}

	return finalParams
}

type ManifestValues struct {
	Values map[string]interface{}
}
