package config

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// ShConfig is a parser loader for .version.sh file
type ShConfig struct {
	values map[string]string
}

// Exists tells if the default location exists for this configType (.version).
func (v *ShConfig) Exists(wd string) bool {
	_, err := os.Stat(filepath.Join(wd, ".version.sh"))
	return !os.IsNotExist(err)
}

// LoadDefault loads the scripts from a default location (.version)
func (v *ShConfig) LoadDefault(wd string) error {
	return v.Load(filepath.Join(wd, ".version.sh"))
}

// GetScript returns the content of the named script.
func (v *ShConfig) GetScript(name string) (string, bool) {
	if x, ok := v.values[name]; ok {
		return x, true
	}
	return "", false
}

// Load given path into the current SimpleConfig object
func (v *ShConfig) Load(path string) error {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}
	return v.Parse(data)
}

// Parse and load given data into the current ShConfig object
func (v *ShConfig) Parse(data []byte) error {
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
	scriptID := regexp.MustCompile(`(?i)^\s*(prebump|prepatch|preminor|premajor|preversion|postversion|postmajor|postminor|postpatch|postbump)=\s*$`)

	lines := lineEnding.Split(string(data), -1)

	currentScriptID := ""
	currentScript := ""
	for _, line := range lines {
		if !comments.MatchString(line) {
			if scriptID.MatchString(line) {
				if currentScriptID != "" {
					v.values[currentScriptID] = strings.TrimSpace(currentScript)
				}
				parts := scriptID.FindAllStringSubmatch(line, -1)
				currentScriptID = strings.ToLower(parts[0][1])
				currentScript = ""
			} else if currentScriptID != "" && line != "" {
				currentScript += strings.TrimLeft(strings.TrimLeft(line, "\t"), " ") + "\n"
			}
		}
	}
	if currentScriptID != "" {
		v.values[currentScriptID] = strings.TrimSpace(currentScript)
	}

	return nil
}
