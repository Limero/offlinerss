package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/limero/offlinerss/domain"
	"github.com/limero/offlinerss/util"
)

func getConfig() (*domain.Config, error) {
	var config domain.Config
	configPath := util.ConfigDir("offlinerss/config.json")
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
			util.WriteFile(string(configJson), configPath)
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
