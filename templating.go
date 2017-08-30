/*
* @Author: Jim Weber
* @Date:   2017-04-23 23:49:30
* @Last Modified by:   Jim Weber
* @Last Modified time: 2017-04-23 23:49:45
 */

package main

import (
	"bytes"
	"log"
	"text/template"
)

func parseManifestTmpl(params ManifestValues, manifestTmpl string) string {

	// t := template.New("manifest-template") //create a new template with some name
	t := template.Must(template.ParseGlob("templates/*.tmpl"))
	_, err := t.Parse(manifestTmpl) //parse some content and generate a template, which is an internal representation
	if err != nil {
		log.Println("Error parsing the  specified template:", err)
	}

	parsedBuffer := new(bytes.Buffer)

	err = t.Execute(parsedBuffer, params)
	if err != nil {
		log.Println("Error executing template:", err)
	}

	return parsedBuffer.String()

}
