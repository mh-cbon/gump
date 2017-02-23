package config

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"gopkg.in/yaml.v2"
)

// GlideConfig can read and load a glide.yaml script node.
type GlideConfig struct {
	Scripts Script `yaml:"scripts"`
}

// Exists tells if the default location exists for this configType (glide.yaml).
func (g *GlideConfig) Exists(wd string) bool {
	_, err := os.Stat(filepath.Join(wd, "glide.yaml"))
	return !os.IsNotExist(err)
}

// LoadDefault loads the scripts from a default location (glide.yaml)
func (g *GlideConfig) LoadDefault(wd string) error {
	return g.Load(filepath.Join(wd, "glide.yaml"))
}

// Script is a struct for the yaml decoding.
type Script struct {
	PreBump     string `yaml:"prebump"`
	PrePatch    string `yaml:"prepatch"`
	PreMinor    string `yaml:"preminor"`
	PreMajor    string `yaml:"premajor"`
	PreVersion  string `yaml:"preversion"`
	PostVersion string `yaml:"postversion"`
	PostMajor   string `yaml:"postmajor"`
	PostMinor   string `yaml:"postminor"`
	PostPatch   string `yaml:"postpatch"`
	PostBump    string `yaml:"postbump"`
}

// Load given path into the current GlideConfig object
func (g *GlideConfig) Load(path string) error {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}
	return g.Parse(data)
}

// Parse and load given data into the current GlideConfig object
func (g *GlideConfig) Parse(data []byte) error {
	return yaml.Unmarshal(data, &g)
}

func parseScript(s string) string {
	lineContinue := regexp.MustCompile(`[\\]\n*`)
	s = lineContinue.ReplaceAllString(s, "")
	return strings.TrimRight(strings.TrimRight(s, " "), "\n")
}

// GetScript returns the content of the named script.
func (g *GlideConfig) GetScript(name string) (string, bool) {
	switch name {
	case "prebump":
		return parseScript(g.Scripts.PreBump), true
	case "prepatch":
		return parseScript(g.Scripts.PrePatch), true
	case "preminor":
		return parseScript(g.Scripts.PreMinor), true
	case "premajor":
		return parseScript(g.Scripts.PreMajor), true
	case "preversion":
		return parseScript(g.Scripts.PreVersion), true
	case "postversion":
		return parseScript(g.Scripts.PostVersion), true
	case "postmajor":
		return parseScript(g.Scripts.PostMajor), true
	case "postminor":
		return parseScript(g.Scripts.PostMinor), true
	case "postpatch":
		return parseScript(g.Scripts.PostPatch), true
	case "postbump":
		return parseScript(g.Scripts.PostBump), true
	}
	return "", false
}
