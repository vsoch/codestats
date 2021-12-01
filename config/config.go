package config

import (
	//	"github.com/mitchellh/mapstructure"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"log"
)

type Conf struct {
	Stats []string `yaml:stats`
}

// read the config and return a config type
func readConfig(yamlStr []byte) Conf {

	// First unmarshall into generic structure
	//var data map[string]interface{}
	// A config can hold multiple keyed sections
	c := Conf{}

	err := yaml.Unmarshal(yamlStr, &c)
	if err != nil {
		log.Fatalf("Unmarshal: %v\n", err)
	}

	// If we have a dockerhierarchy, add it
	//if item, ok := data["dockerhierarchy"]; ok {
	//	c.DockerHierarchy = convertDockerHierarchy(item)
	//}

	return c
}

// convertDockerHierarchy maps the dockerhierarchy portion to a DockerHierarchy
//func convertDockerHierarchy(item interface{}) DockerHierarchy {
//	hier := DockerHierarchy{}
//	mapstructure.Decode(item, &hier)
//	return hier
//}

func Load(yamlfile string) Conf {
	yamlContent, err := ioutil.ReadFile(yamlfile)
	if err != nil {
		log.Printf("yamlFile.Get err   #%v ", err)
	}
	return readConfig(yamlContent)
}
