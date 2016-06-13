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
