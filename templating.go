/*
* @Author: Jim Weber
* @Date:   2017-04-23 23:49:30
* @Last Modified by:   Jim Weber
* @Last Modified time: 2017-04-23 23:49:45
 */

package main

import (
	"bytes"
	"text/template"
)

func parseManifestTmpl(params ManifestValues, manifestTmpl string) string {

	t := template.New("manifest template") //create a new template with some name
	t, _ = t.Parse(manifestTmpl)           //parse some content and generate a template, which is an internal representation

	parsedBuffer := new(bytes.Buffer)
	t.Execute(parsedBuffer, params)

	return parsedBuffer.String()

}
