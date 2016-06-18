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
	v.values["preversion"] = ""
	v.values["postversion"] = ""

	lineEnding := regexp.MustCompile("\r\n|\n")
	comments := regexp.MustCompile(`^\s*#`)
	lineContinue := regexp.MustCompile(`[\\]\n*$`)
	scriptId := regexp.MustCompile(`^\s*(preversion|postversion):\s*(.+)`)

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

// GetPreVersion returns the content of preversion field
func (v *SimpleConfig) GetPreVersion() string {
	return v.values["preversion"]
}

// GetPostVersion returns the content of postversion field
func (v *SimpleConfig) GetPostVersion() string {
	return v.values["postversion"]
}
