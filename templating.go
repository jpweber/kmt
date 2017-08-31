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

func initTemplate() *template.Template {
	// start by trying to init with the templates dir
	// to include all subtemplates and such.
	// However this is path specific and doesn't always work
	// come up with a way to fix this that isn't just passing the template path
	// as a CLI arg.
	t, err := template.ParseGlob("templates/*.tmpl")
	if err != nil {
		// log the error in debug mode
		logger(err.Error())

		// init our own empty template if the parseGlob fails
		t = template.New("manifest-template") //create a new template with some name
		return t
	}

	// if the parseglob succeded return that
	return t
}

func parseManifestTmpl(params ManifestValues, manifestTmpl string) string {

	t := initTemplate()

	// going to take advantage of this later.
	// add stringJoin function to templated
	// _ = t.Funcs(template.FuncMap{"StringsJoin": strings.Join})
	// if err != nil {
	// 	log.Println("Error adding function to the specified template:", err)
	// }

	// parse the user specified manifest template
	_, err := t.Parse(manifestTmpl)
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
