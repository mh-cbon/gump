package config

import (
	"fmt"
	"testing"
)

func TestEmptyYamlFile(t *testing.T) {
	s := &GlideConfig{}
	s.Parse([]byte(""))
	expects := map[string]string{}
	expects["preversion"] = ""
	expects["postversion"] = ""
	checkTable(t, expects, s)
}

func TestYamlFile(t *testing.T) {
	s := &GlideConfig{}
	s.Parse([]byte(`
scripts:
  preversion: some
  postversion: else
`))
	expects := map[string]string{}
	expects["preversion"] = "some"
	expects["postversion"] = "else"
	checkTable(t, expects, s)
}

func TestYamlFileMultiline(t *testing.T) {
	s := &GlideConfig{}
	s.Parse([]byte(`
scripts:
  preversion: |
    some \
    thing
  postversion: |
    else \
      otherwise
`))
	expects := map[string]string{}
	expects["preversion"] = "some thing"
	expects["postversion"] = "else   otherwise"
	checkTable(t, expects, s)
}

func TestAllFieldsGlide(t *testing.T) {
	a := []string{
		"prebump",
		"prepatch",
		"preminor",
		"premajor",
		"preversion",
		"postversion",
		"postmajor",
		"postminor",
		"postpatch",
		"postbump",
	}
	scriptStr := "\nscripts:\n"
	for _, phase := range a {
		scriptStr += fmt.Sprintf("  %v: %v\n", phase, phase)
	}
	scriptStr += "\n"

	s := &GlideConfig{}
	s.Parse([]byte(scriptStr))

	expects := map[string]string{}
	for _, phase := range a {
		expects[phase] = phase
	}
	checkTable(t, expects, s)
}
