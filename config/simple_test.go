package config

import (
	"testing"
)

func TestEmptyFile(t *testing.T) {
	s := SimpleConfig{}
	s.Parse([]byte(""))
	expected := ""
	if s.values["preversion"] != expected {
		t.Errorf("Expected %q got %q\n", expected, s.values["preversion"])
	}
	if s.values["postversion"] != expected {
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
	if s.values["preversion"] != expected {
		t.Errorf("Expected %q got %q\n", expected, s.values["preversion"])
	}
	expected = "else"
	if s.values["postversion"] != expected {
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
	if s.values["preversion"] != expected {
		t.Errorf("Expected %q got %q\n", expected, s.values["preversion"])
	}
	expected = "else   otherwise"
	if s.values["postversion"] != expected {
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
	if s.values["preversion"] != expected {
		t.Errorf("Expected %q got %q\n", expected, s.values["preversion"])
	}
	expected = "else"
	if s.values["postversion"] != expected {
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
	if s.values["preversion"] != expected {
		t.Errorf("Expected %q got %q\n", expected, s.values["preversion"])
	}
	expected = "else"
	if s.values["postversion"] != expected {
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
	if s.values["preversion"] != expected {
		t.Errorf("Expected %q got %q\n", expected, s.values["preversion"])
	}
	expected = "else "
	if s.values["postversion"] != expected {
		t.Errorf("Expected %q got %q\n", expected, s.values["postversion"])
	}
}

func TestAllFieldsSimpleConfig(t *testing.T) {
	s := SimpleConfig{}
	s.Parse([]byte(`
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
