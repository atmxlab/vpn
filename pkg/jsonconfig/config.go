package jsonconfig

import (
	"encoding/json"
	"os"
)

func Load[T any](path string) (T, error) {
	var config T
	file, err := os.Open(path)
	if err != nil {
		return config, err
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	err = decoder.Decode(&config)
	if err != nil {
		return config, err
	}

	return config, nil
}
