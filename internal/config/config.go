package config

import (
	"encoding/json"
	"io/ioutil"
)

func ReadConfig() (*Data, error) {
	fileContent, err := ioutil.ReadFile("data/input/config.json")
	if err != nil {
		return nil, err
	}

	var data Data
	err = json.Unmarshal(fileContent, &data)
	if err != nil {
		return nil, err
	}

	return &data, nil
}
