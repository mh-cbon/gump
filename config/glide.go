package config

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"regexp"
	"strings"
)

// Config is the top-level configuration object.
type GlideConfig struct {
	Scripts Script `yaml:"scripts"`
}

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

// GetPreBump returns the content of scripts.prebump field
func (g *GlideConfig) GetPreBump() string {
	return parseScript(g.Scripts.PreBump)
}

// GetPrePatch returns the content of scripts.prepatch field
func (g *GlideConfig) GetPrePatch() string {
	return parseScript(g.Scripts.PrePatch)
}

// GetPreMinor returns the content of scripts.preminor field
func (g *GlideConfig) GetPreMinor() string {
	return parseScript(g.Scripts.PreMinor)
}

// GetPreMajor returns the content of scripts.premajor field
func (g *GlideConfig) GetPreMajor() string {
	return parseScript(g.Scripts.PreMajor)
}

// GetPreVersion returns the content of scripts.preversion field
func (g *GlideConfig) GetPreVersion() string {
	return parseScript(g.Scripts.PreVersion)
}

// GetPostVersion returns the content of scripts.postversion field
func (g *GlideConfig) GetPostVersion() string {
	return parseScript(g.Scripts.PostVersion)
}

// GetPostMajor returns the content of scripts.postmajor field
func (g *GlideConfig) GetPostMajor() string {
	return parseScript(g.Scripts.PostMajor)
}

// GetPostMinor returns the content of scripts.postminor field
func (g *GlideConfig) GetPostMinor() string {
	return parseScript(g.Scripts.PostMinor)
}

// GetPostPatch returns the content of scripts.postpatch field
func (g *GlideConfig) GetPostPatch() string {
	return parseScript(g.Scripts.PostPatch)
}

// GetPostBump returns the content of scripts.postbump field
func (g *GlideConfig) GetPostBump() string {
	return parseScript(g.Scripts.PostBump)
}
