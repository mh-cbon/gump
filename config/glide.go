package config

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
  "strings"
  "regexp"
)

// Config is the top-level configuration object.
type GlideConfig struct {
	Scripts Script `yaml:"scripts"`
}

type Script struct {
	PreVersion  string `yaml:"preversion"`
	PostVersion string `yaml:"postversion"`
}

func (g *GlideConfig) Load(path string) error {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}
	return g.Parse(data)
}

func (g *GlideConfig) Parse(data []byte) error {
	return yaml.Unmarshal(data, &g)
}

func (g *GlideConfig) GetPreVersion() string {
  s := g.Scripts.PreVersion
	lineContinue := regexp.MustCompile(`[\\]\n*`)
  s = lineContinue.ReplaceAllString(s, "")
	return strings.TrimRight(strings.TrimRight(s, " "), "\n")
}

func (g *GlideConfig) GetPostVersion() string {
  s := g.Scripts.PostVersion
	lineContinue := regexp.MustCompile(`[\\]\n*`)
  s = lineContinue.ReplaceAllString(s, "")
	return strings.TrimRight(strings.TrimRight(s, " "), "\n")
}
