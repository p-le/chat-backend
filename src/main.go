package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"gopkg.in/yaml.v2"
)

type configuration struct {
	Enabled bool
	Path    string
}

func main() {
	file, _ := os.Open("config.json")
	defer file.Close()
	decoder := json.NewDecoder(file)
	jsonConfig := configuration{}
	err := decoder.Decode(&jsonConfig)
	if err != nil {
		fmt.Println("Error: ", err)
	}
	fmt.Println(jsonConfig.Path)

	yamlConfig := configuration{}
	bytes, _ := ioutil.ReadFile("config.yml")
	error := yaml.Unmarshal(bytes, &yamlConfig)
	if error != nil {
		log.Fatalf("Error: %v", err)
	}
	fmt.Println(yamlConfig.Path)
}
