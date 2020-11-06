package main

import (
	"io/ioutil"
	"log"
	"gopkg.in/yaml.v3"
)

func Find(slice []string, val string) bool {
	for _, item := range slice {
		if item == val {
			return true
		}
	}
	return false
}

func (p *KeyvaultProps) getConf(EventGridFile *string) *KeyvaultProps {

	yamlFile, err := ioutil.ReadFile(*EventGridFile)
	if err != nil {
		log.Fatalf("yamlFile.Get err   #%v ", err)
	}
	err = yaml.Unmarshal(yamlFile, p)
	if err != nil {
		log.Fatalf("Unmarshal: %v", err)
	}

	return p
}
