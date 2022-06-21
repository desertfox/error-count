package storage

import (
	"io/ioutil"
	"log"

	"gopkg.in/yaml.v3"
)

func Load(path string) map[string]int {
	var totals map[string]int = make(map[string]int)

	b, err := ioutil.ReadFile(path)
	if err != nil {
		log.Println(err)
		return totals
	}

	err = yaml.Unmarshal(b, &totals)
	if err != nil {
		log.Println(err)
		return totals
	}

	return totals
}

func Save(path string, totals map[string]int) error {
	data, err := yaml.Marshal(totals)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(path, data, 0644)
}
