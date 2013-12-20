package main

import (
	"encoding/json"
	"fmt"
	"os"
)

type Config struct {
	Db struct {
		Dsn map[string]string
	}

	Html struct {
		HtmlSrcDir string
	}
}

func loadConfig(cfgFileName string) (*Config, error) {
	cfgFile, err := os.Open(cfgFileName)
	if err != nil {
		return nil, fmt.Errorf("Error opening config file %s: %v", cfgFileName, err)
	}
	defer cfgFile.Close()
	cfg := &Config{}
	decoder := json.NewDecoder(cfgFile)
	err = decoder.Decode(cfg)
	if err != nil {
		return nil, fmt.Errorf("Error parsing config file: %v", err)
	}
	return cfg, nil
}
