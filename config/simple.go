package config

import (
	"io/ioutil"
	"regexp"
)

// SimpleConfig is a parser loader for .version file
type SimpleConfig struct {
	values map[string]string
}

// Load given path into the current SimpleConfig object
func (v *SimpleConfig) Load(path string) error {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}
	return v.Parse(data)
}

// Parse and load given data into the current SimpleConfig object
func (v *SimpleConfig) Parse(data []byte) error {
	v.values = make(map[string]string)
	v.values["prebump"] = ""
	v.values["prepatch"] = ""
	v.values["preminor"] = ""
	v.values["premajor"] = ""
	v.values["preversion"] = ""
	v.values["postversion"] = ""
	v.values["postmajor"] = ""
	v.values["postminor"] = ""
	v.values["postpatch"] = ""
	v.values["postbump"] = ""

	lineEnding := regexp.MustCompile("\r\n|\n")
	comments := regexp.MustCompile(`^\s*#`)
	lineContinue := regexp.MustCompile(`[\\]\n*$`)
	scriptId := regexp.MustCompile(`^\s*(prebump|prepatch|preminor|premajor|preversion|postversion|postmajor|postminor|postpatch|postbump):\s*(.+)`)

	lines := lineEnding.Split(string(data), -1)

	isContinuing := false
	currentScriptId := ""
	currentScript := ""
	for _, line := range lines {
		if comments.MatchString(line) {
		} else if scriptId.MatchString(line) {
			if currentScriptId != "" {
				v.values[currentScriptId] = currentScript
			}
			isContinuing = lineContinue.MatchString(line)
			parts := scriptId.FindAllStringSubmatch(line, -1)
			currentScriptId = parts[0][1]
			currentScript = parts[0][2]
			if isContinuing {
				currentScript = lineContinue.ReplaceAllString(currentScript, "")
			}
		} else if isContinuing {
			currentScript += lineContinue.ReplaceAllString(line, "")
			isContinuing = lineContinue.MatchString(line)
		}
	}
	if currentScriptId != "" {
		v.values[currentScriptId] = currentScript
	}

	return nil
}

// GetPreBump returns the content of prebump field
func (v *SimpleConfig) GetPreBump() string {
	return v.values["prebump"]
}

// GetPrePatch returns the content of prepatch field
func (v *SimpleConfig) GetPrePatch() string {
	return v.values["prepatch"]
}

// GetPreMinor returns the content of preminor field
func (v *SimpleConfig) GetPreMinor() string {
	return v.values["preminor"]
}

// GetPreMajor returns the content of premajor field
func (v *SimpleConfig) GetPreMajor() string {
	return v.values["premajor"]
}

// GetPreVersion returns the content of preversion field
func (v *SimpleConfig) GetPreVersion() string {
	return v.values["preversion"]
}

// GetPostVersion returns the content of postversion field
func (v *SimpleConfig) GetPostVersion() string {
	return v.values["postversion"]
}

// GetPostMajor returns the content of postmajor field
func (v *SimpleConfig) GetPostMajor() string {
	return v.values["postmajor"]
}

// GetPostMinor returns the content of postminor field
func (v *SimpleConfig) GetPostMinor() string {
	return v.values["postminor"]
}

// GetPostPatch returns the content of postpatch field
func (v *SimpleConfig) GetPostPatch() string {
	return v.values["postpatch"]
}

// GetPostBump returns the content of postbump field
func (v *SimpleConfig) GetPostBump() string {
	return v.values["postbump"]
}
