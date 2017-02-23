package stringexec

import (
	"io/ioutil"
	"os"
	"os/exec"
	"runtime"
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

// NewTempCmd ...
func NewTempCmd(cwd string, cmd string) (*TempCmd, error) {
	f, err := ioutil.TempDir("", "stringexec")
	if err != nil {
		return nil, err
	}
	fp := f + "/s"
	err = ioutil.WriteFile(fp, []byte(cmd), 0766)
	if err != nil {
		return nil, err
	}
	ret := &TempCmd{Cmd: exec.Command("sh", "-c", fp), f: fp}
	if runtime.GOOS == "windows" {
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
