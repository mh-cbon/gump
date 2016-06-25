package config

import (
	"testing"
)

func TestEmptyYamlFile(t *testing.T) {
	s := GlideConfig{}
	s.Parse([]byte(""))
	expected := ""
	if s.GetPreVersion() != expected {
		t.Errorf("Expected %q got %q\n", expected, s.GetPreVersion())
	}
	if s.GetPostVersion() != expected {
		t.Errorf("Expected %q got %q\n", expected, s.GetPostVersion())
	}
}

func TestYamlFile(t *testing.T) {
	s := GlideConfig{}
	s.Parse([]byte(`
scripts:
  preversion: some
  postversion: else
`))
	expected := "some"
	if s.GetPreVersion() != expected {
		t.Errorf("Expected %q got %q\n", expected, s.GetPreVersion())
	}
	expected = "else"
	if s.GetPostVersion() != expected {
		t.Errorf("Expected %q got %q\n", expected, s.GetPostVersion())
	}
}

func TestYamlFileMultiline(t *testing.T) {
	s := GlideConfig{}
	s.Parse([]byte(`
scripts:
  preversion: |
    some \
    thing
  postversion: |
    else \
      otherwise
`))
	expected := "some thing"
	if s.GetPreVersion() != expected {
		t.Errorf("Expected %q got %q\n", expected, s.GetPreVersion())
	}
	expected = "else   otherwise"
	if s.GetPostVersion() != expected {
		t.Errorf("Expected %q got %q\n", expected, s.GetPostVersion())
	}
}

func TestAllFields(t *testing.T) {
	s := GlideConfig{}
	s.Parse([]byte(`
scripts:
  prebump: prebump
  prepatch: prepatch
  preminor: preminor
  premajor: premajor
  preversion: preversion
  postversion: postversion
  postmajor: postmajor
  postminor: postminor
  postpatch: postpatch
  postbump: postbump
`))
	expected := "prebump"
	if s.GetPreBump() != expected {
		t.Errorf("Expected %q got %q\n", expected, s.GetPreBump())
	}
	expected = "prepatch"
	if s.GetPrePatch() != expected {
		t.Errorf("Expected %q got %q\n", expected, s.GetPrePatch())
	}
	expected = "preminor"
	if s.GetPreMinor() != expected {
		t.Errorf("Expected %q got %q\n", expected, s.GetPreMinor())
	}
	expected = "premajor"
	if s.GetPreMajor() != expected {
		t.Errorf("Expected %q got %q\n", expected, s.GetPreMajor())
	}
	expected = "preversion"
	if s.GetPreVersion() != expected {
		t.Errorf("Expected %q got %q\n", expected, s.GetPreVersion())
	}
	expected = "postversion"
	if s.GetPostVersion() != expected {
		t.Errorf("Expected %q got %q\n", expected, s.GetPostVersion())
	}
	expected = "postmajor"
	if s.GetPostMajor() != expected {
		t.Errorf("Expected %q got %q\n", expected, s.GetPostMajor())
	}
	expected = "postminor"
	if s.GetPostMinor() != expected {
		t.Errorf("Expected %q got %q\n", expected, s.GetPostMinor())
	}
	expected = "postpatch"
	if s.GetPostPatch() != expected {
		t.Errorf("Expected %q got %q\n", expected, s.GetPostPatch())
	}
	expected = "postbump"
	if s.GetPostBump() != expected {
		t.Errorf("Expected %q got %q\n", expected, s.GetPostBump())
	}
}
