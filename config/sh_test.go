package config

import (
	"fmt"
	"testing"
)

func TestEmptyShFile(t *testing.T) {
	s := &ShConfig{}
	s.Parse([]byte(""))
	expects := map[string]string{}
	expects["preversion"] = ""
	expects["postversion"] = ""
	checkTable(t, expects, s)
}

func TestShFile(t *testing.T) {
	s := &ShConfig{}
	s.Parse([]byte(`
preversion=
  some

postversion=
  else
`))
	expects := map[string]string{}
	expects["preversion"] = "some"
	expects["postversion"] = "else"
	checkTable(t, expects, s)
}

func TestShFileMultiline(t *testing.T) {
	s := &ShConfig{}
	s.Parse([]byte(`
preversion=
  some \
  thing

postversion=
  else \
    otherwise
`))
	expects := map[string]string{}
	expects["preversion"] = "some \\\nthing"
	expects["postversion"] = "else \\\notherwise"
	checkTable(t, expects, s)
}

func TestAllFieldsSh(t *testing.T) {
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
		scriptStr += fmt.Sprintf("%v=\n", phase)
		scriptStr += fmt.Sprintf("\t%v\n", phase)
	}
	scriptStr += "\n"

	s := &ShConfig{}
	s.Parse([]byte(scriptStr))

	expects := map[string]string{}
	for _, phase := range a {
		expects[phase] = phase
	}
	checkTable(t, expects, s)
}

func checkTable(t *testing.T, expects map[string]string, l Configured) bool {
	for phase, expect := range expects {
		if got, found := l.GetScript(phase); !found || got != expect {
			t.Errorf("loader(%T) phase(%v) expected=%q got=%q\n", l, phase, expect, got)
			return false
		}
	}
	return true
}
