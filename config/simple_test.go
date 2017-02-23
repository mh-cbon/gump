package config

import (
	"fmt"
	"testing"
)

func TestEmptyFile(t *testing.T) {
	s := &SimpleConfig{}
	s.Parse([]byte(""))
	expects := map[string]string{}
	expects["preversion"] = ""
	expects["postversion"] = ""
	checkTable(t, expects, s)
}

func TestFile(t *testing.T) {
	s := &SimpleConfig{}
	s.Parse([]byte(`
preversion: some
postversion: else
`))
	expects := map[string]string{}
	expects["preversion"] = "some"
	expects["postversion"] = "else"
	checkTable(t, expects, s)
}

func TestFileMultiline(t *testing.T) {
	s := &SimpleConfig{}
	s.Parse([]byte(`
preversion: some \
thing
postversion: else \
  otherwise
`))
	expects := map[string]string{}
	expects["preversion"] = "some thing"
	expects["postversion"] = "else   otherwise"
	checkTable(t, expects, s)
}

func TestFileWithComments(t *testing.T) {
	s := &SimpleConfig{}
	s.Parse([]byte(`
# comment1
preversion: some
# comment1
postversion: else
# comment2
`))
	expects := map[string]string{}
	expects["preversion"] = "some"
	expects["postversion"] = "else"
	checkTable(t, expects, s)
}

func TestFileMultilineWithComments(t *testing.T) {
	s := &SimpleConfig{}
	s.Parse([]byte(`
preversion: some \
# comment1
thing
postversion: else
`))
	expects := map[string]string{}
	expects["preversion"] = "some thing"
	expects["postversion"] = "else"
	checkTable(t, expects, s)
}

func TestFileMalformedMultiline(t *testing.T) {
	s := &SimpleConfig{}
	s.Parse([]byte(`
preversion: some \
postversion: else \
`))
	expects := map[string]string{}
	expects["preversion"] = "some "
	expects["postversion"] = "else "
	checkTable(t, expects, s)
}

func TestAllFieldsSimpleConfig(t *testing.T) {
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
	scriptStr := ""
	for _, phase := range a {
		scriptStr += fmt.Sprintf("%v: %v\n", phase, phase)
	}

	s := &SimpleConfig{}
	s.Parse([]byte(scriptStr))

	expects := map[string]string{}
	for _, phase := range a {
		expects[phase] = phase
	}
	checkTable(t, expects, s)
}
