package config

import (
	"errors"
	"os"
)

type Configured interface {
	Load(path string) error
	GetPreVersion() string
	GetPostVersion() string
}

func Exists(path string) bool {
	if _, err := os.Stat(path + "/.version"); !os.IsNotExist(err) {
		return true
	}
	if _, err := os.Stat(path + "/glide.yaml"); !os.IsNotExist(err) {
		return true
	}
	return false
}

func Load(path string) (Configured, error) {
	var config Configured
	configType := ""
	filepath := path + "/.version"
	if _, err := os.Stat(filepath); !os.IsNotExist(err) {
		configType = ".version"
		filepath = path + "/.version"
	} else {
		filepath = path + "/glide.yaml"
		if _, err = os.Stat(filepath); !os.IsNotExist(err) {
			configType = "glide"
		}
	}
	if configType == "" {
		return nil, errors.New("Cannot find suitable file to load the version scripts.")
	} else if configType == ".version" {
		config = &SimpleConfig{}
	} else if configType == "glide" {
		config = &GlideConfig{}
	}
	err := config.Load(filepath)
	return config, err
}
