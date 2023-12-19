package config

import (
	"encoding/json"
	"os"
)

func ParseFile(cfg any, filePath string) error {
	f, err := os.Open(filePath)
	if err != nil {
		return err
	}

	return json.NewDecoder(f).Decode(cfg)
}
