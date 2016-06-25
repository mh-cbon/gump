package config

import (
	"errors"
	"os"
)

// Configured is a facade to a glide.yaml / . version script file
type Configured interface {
	Load(path string) error
	GetPreBump() string
	GetPrePatch() string
	GetPreMinor() string
	GetPreMajor() string
	GetPreVersion() string
	GetPostVersion() string
	GetPostMajor() string
	GetPostMinor() string
	GetPostPatch() string
	GetPostBump() string
}

// Exists tells if a version/glide.yaml file exists at given directory
func Exists(path string) bool {
	if _, err := os.Stat(path + "/.version"); !os.IsNotExist(err) {
		return true
	}
	if _, err := os.Stat(path + "/glide.yaml"); !os.IsNotExist(err) {
		return true
	}
	return false
}

// Load version script from the given directory
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

// get a script given its string name
func GetScript(which string, c Configured) string {
	if which == "prebump" {
		return c.GetPreBump()

	} else if which == "prepatch" {
		return c.GetPrePatch()

	} else if which == "preminor" {
		return c.GetPreMinor()

	} else if which == "premajor" {
		return c.GetPreMajor()

	} else if which == "preversion" {
		return c.GetPreVersion()

	} else if which == "postversion" {
		return c.GetPostVersion()

	} else if which == "postmajor" {
		return c.GetPostMajor()

	} else if which == "postminor" {
		return c.GetPostMinor()

	} else if which == "postpatch" {
		return c.GetPostPatch()

	} else if which == "postbump" {
		return c.GetPostBump()

	}
	return ""
}
