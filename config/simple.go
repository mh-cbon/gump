package config

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
)

// SimpleConfig is a parser loader for .version file
type SimpleConfig struct {
	values map[string]string
}

// Exists tells if the default location exists for this configType (.version).
func (v *SimpleConfig) Exists(wd string) bool {
	_, err := os.Stat(filepath.Join(wd, ".version"))
	return !os.IsNotExist(err)
}

// LoadDefault loads the scripts from a default location (.version)
func (v *SimpleConfig) LoadDefault(wd string) error {
	return v.Load(filepath.Join(wd, ".version"))
}

// GetScript returns the content of the named script.
func (v *SimpleConfig) GetScript(name string) (string, bool) {
	if x, ok := v.values[name]; ok {
		return x, true
	}
	return "", false
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
	scriptID := regexp.MustCompile(`^\s*(prebump|prepatch|preminor|premajor|preversion|postversion|postmajor|postminor|postpatch|postbump):\s*(.+)`)

	lines := lineEnding.Split(string(data), -1)

	isContinuing := false
	currentScriptID := ""
	currentScript := ""
	for _, line := range lines {
		if comments.MatchString(line) {
		} else if scriptID.MatchString(line) {
			if currentScriptID != "" {
				v.values[currentScriptID] = currentScript
			}
			isContinuing = lineContinue.MatchString(line)
			parts := scriptID.FindAllStringSubmatch(line, -1)
			currentScriptID = parts[0][1]
			currentScript = parts[0][2]
			if isContinuing {
				currentScript = lineContinue.ReplaceAllString(currentScript, "")
			}
		} else if isContinuing {
			currentScript += lineContinue.ReplaceAllString(line, "")
			isContinuing = lineContinue.MatchString(line)
		}
	}
	if currentScriptID != "" {
		v.values[currentScriptID] = currentScript
	}

	return nil
}
