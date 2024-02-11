package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/limero/offlinerss/helpers"
	"github.com/limero/offlinerss/models"
)

func getConfig() (*models.Config, error) {
	var config models.Config
	configPath := filepath.Join(helpers.ConfigDir(), "offlinerss/config.json")
	data, err := os.ReadFile(configPath)
	if err != nil {
		if os.IsNotExist(err) {
			config, err = setup()
			if err != nil {
				return nil, err
			}

			configJson, err := json.MarshalIndent(config, "", "  ")
			if err != nil {
				return nil, err
			}
			helpers.WriteFile(string(configJson), configPath)
			fmt.Printf("Successfully written new config to %q\n\n", configPath)

			return &config, nil
		}
		return nil, err
	}

	if err = json.Unmarshal(data, &config); err != nil {
		return nil, err
	}

	return &config, nil
}
