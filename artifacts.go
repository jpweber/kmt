package main

import (
	"bytes"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"sync"
)

func writeArtifacts(manPath, noExtFileName, manifest string, parameters ManifestValues) {
	// by default write out the file to the following path
	// ./artifacts/<realm>/<namespace>/<app name>/filename
	// a dir container latest versions of the file is also use
	// ./artifacts/<realm>/<namespace>/current/filename
	// if a realm parameter does not exist it will be exlucded from the path.
	// if a name space is not used "default" will be put in the path
	// place artifacts file in the root dir of the manifests
	var namespace string

	// handle blank envivironments
	if parameters.Values["environ"] == nil {
		namespace = "default"
	} else {
		namespace = parameters.Values["environ"].(string)
	}

	rootArt, _ := filepath.Split(manPath)
	var artPath bytes.Buffer
	artPath.WriteString(rootArt)
	artPath.WriteString("/")
	artPath.WriteString("artifacts")
	artPath.WriteString("/")
	// include the realm in the path only if it is provided
	if parameters.Values["realm"] != nil {
		artPath.WriteString(parameters.Values["realm"].(string))
		artPath.WriteString("/")
	}
	artPath.WriteString(namespace)
	artPath.WriteString("/")
	// get the path for _current_ files before completing the path
	latestArtPath := artPath.String()

	artPath.WriteString(parameters.Values["name"].(string))
	artPath.WriteString("/")
	err := os.MkdirAll(artPath.String(), 0755)
	if err != nil {
		log.Println("Error making dir")
	}

	// make the environment specific latest dir
	err = os.MkdirAll(latestArtPath+"/latest", 0755)
	if err != nil {
		log.Println("Error making latest dir")
	}

	// versioned artifact filename
	var artFile bytes.Buffer
	artFile.WriteString(artPath.String())
	artFile.WriteString("/")
	artFile.WriteString(parameters.Values["name"].(string) + "-" + parameters.Values["version"].(string) + ".yaml")

	// latest artifact filename
	var latestArtFile bytes.Buffer
	latestArtFile.WriteString(latestArtPath)
	latestArtFile.WriteString("/latest/")
	latestArtFile.WriteString(parameters.Values["name"].(string) + ".yaml")

	var wg sync.WaitGroup

	wg.Add(2)
	// write the data in to the file
	go func(wg *sync.WaitGroup) {
		err = ioutil.WriteFile(artFile.String(), []byte(manifest), 0755)
		if err != nil {
			log.Println("There was an error writing your manifest artifact:", err)

		}
		wg.Done()
	}(&wg)

	// write the current file
	go func(wg *sync.WaitGroup) {
		err = ioutil.WriteFile(latestArtFile.String(), []byte(manifest), 0755)
		if err != nil {
			log.Println("There was an error writing your _current_ manifest artifact:", err)

		}
		wg.Done()
	}(&wg)

	wg.Wait()
}
