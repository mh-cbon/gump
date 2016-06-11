package config

import (
	"testing"
)

func TestEmptyFile(t *testing.T) {
  s := SimpleConfig{}
  s.Parse([]byte(""))
  expected := ""
	if s.values["preversion"]!=expected {
		t.Errorf("Expected %q got %q\n", expected, s.values["preversion"])
	}
	if s.values["postversion"]!=expected {
		t.Errorf("Expected %q got %q\n", expected, s.values["postversion"])
	}
}

func TestFile(t *testing.T) {
  s := SimpleConfig{}
  s.Parse([]byte(`
preversion: some
postversion: else
`))
  expected := "some"
	if s.values["preversion"]!=expected {
		t.Errorf("Expected %q got %q\n", expected, s.values["preversion"])
	}
  expected = "else"
	if s.values["postversion"]!=expected {
		t.Errorf("Expected %q got %q\n", expected, s.values["postversion"])
	}
}

func TestFileMultiline(t *testing.T) {
  s := SimpleConfig{}
  s.Parse([]byte(`
preversion: some \
thing
postversion: else \
  otherwise
`))
  expected := "some thing"
	if s.values["preversion"]!=expected {
		t.Errorf("Expected %q got %q\n", expected, s.values["preversion"])
	}
  expected = "else   otherwise"
	if s.values["postversion"]!=expected {
		t.Errorf("Expected %q got %q\n", expected, s.values["postversion"])
	}
}

func TestFileWithComments(t *testing.T) {
  s := SimpleConfig{}
  s.Parse([]byte(`
# comment1
preversion: some
# comment1
postversion: else
# comment2
`))
  expected := "some"
	if s.values["preversion"]!=expected {
		t.Errorf("Expected %q got %q\n", expected, s.values["preversion"])
	}
  expected = "else"
	if s.values["postversion"]!=expected {
		t.Errorf("Expected %q got %q\n", expected, s.values["postversion"])
	}
}

func TestFileMultilineWithComments(t *testing.T) {
  s := SimpleConfig{}
  s.Parse([]byte(`
preversion: some \
# comment1
thing
postversion: else
`))
  expected := "some thing"
	if s.values["preversion"]!=expected {
		t.Errorf("Expected %q got %q\n", expected, s.values["preversion"])
	}
  expected = "else"
	if s.values["postversion"]!=expected {
		t.Errorf("Expected %q got %q\n", expected, s.values["postversion"])
	}
}

func TestFileMalformedMultiline(t *testing.T) {
  s := SimpleConfig{}
  s.Parse([]byte(`
preversion: some \
postversion: else \
`))
  expected := "some "
	if s.values["preversion"]!=expected {
		t.Errorf("Expected %q got %q\n", expected, s.values["preversion"])
	}
  expected = "else "
	if s.values["postversion"]!=expected {
		t.Errorf("Expected %q got %q\n", expected, s.values["postversion"])
	}
}
