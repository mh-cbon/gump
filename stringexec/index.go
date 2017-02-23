package stringexec

import (
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
)

// Command Return a new exec.Cmd object for the given command string
func Command(cwd string, cmd string) (*TempCmd, error) {
	return NewTempCmd(cwd, cmd)
}

// TempCmd ...
type TempCmd struct {
	*exec.Cmd
	f string
}

var isWindows = runtime.GOOS == "windows"

func prepareCommand(cmd string, isWindows bool) string {
	ret := ""
	lineEnding := regexp.MustCompile("\r\n|\n")
	continueStart := regexp.MustCompile(`^\s*&&\s*`)
	continueEnd := regexp.MustCompile(`[\\]\s*$`)
	lines := lineEnding.Split(cmd, -1)

	if isWindows {
		for _, line := range lines {
			if continueEnd.MatchString(line) {
				line = continueEnd.ReplaceAllString(line, "")
			} else {
				line += " && "
			}
			ret += strings.TrimSpace(line) + " "
		}
		ret = strings.TrimSpace(ret)
		if strings.HasSuffix(ret, " &&") {
			ret = ret[0 : len(ret)-3]
		}
	} else {
		isContinuing := false
		for i, line := range lines {
			if i == 0 {
				isContinuing = false
				if !continueEnd.MatchString(line) {
					line += " \\"
					isContinuing = true
				}
			} else {
				if isContinuing && !continueStart.MatchString(line) {
					line = " && " + line
				}
				isContinuing = false
				if !continueEnd.MatchString(line) {
					line += " \\"
					isContinuing = true
				}
			}
			line += "\n"
			ret += line
		}
		if len(ret) > 3 {
			ret = ret[0 : len(ret)-3]
		}
	}
	return ret
}

// NewTempCmd ...
func NewTempCmd(cwd string, cmd string) (*TempCmd, error) {
	f, err := ioutil.TempDir("", "stringexec")
	if err != nil {
		return nil, err
	}
	fp := filepath.Join(f, "s")
	if isWindows {
		fp += ".bat"
	}
	cmd = prepareCommand(cmd, isWindows)
	err = ioutil.WriteFile(fp, []byte(cmd), 0766)
	if err != nil {
		return nil, err
	}
	ret := &TempCmd{Cmd: exec.Command("sh", "-c", fp), f: fp}
	if isWindows {
		ret.Cmd = exec.Command("cmd", "/C", fp)
	}
	ret.Cmd.Dir = cwd
	ret.Cmd.Stdout = os.Stdout
	ret.Cmd.Stderr = os.Stderr
	return ret, nil
}

// Run ...
func (t *TempCmd) Run() error {
	if err := t.Cmd.Start(); err != nil {
		return err
	}
	return t.Wait()
}

// Wait ...
func (t *TempCmd) Wait() error {
	err := t.Cmd.Wait()
	os.Remove(t.f)
	return err
}
