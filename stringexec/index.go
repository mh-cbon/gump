package stringexec

import (
	"io/ioutil"
	"os/exec"
	"runtime"
)

func Command(cmd string) (*exec.Cmd, error) {
	if runtime.GOOS == "windows" {
		return ExecStringWindows(cmd)
	}
	return ExecStringFriendlyUnix(cmd)
}

func ExecStringWindows(cmd string) (*exec.Cmd, error) {
	dir, err := ioutil.TempDir("", "stringexec")
	if err != nil {
		return nil, err
	}
	err = ioutil.WriteFile(dir+"/some.bat", []byte(cmd), 0766)
	if err != nil {
		return nil, err
	}

	oCmd := exec.Command("cmd", []string{dir + "/some.bat"}...)
  oCmd.Dir = dir
	// defer os.Remove(tmpfile.Name()) // clean up // not sure how to clean it :x
	return oCmd, nil
}

func ExecStringFriendlyUnix(cmd string) (*exec.Cmd, error) {
	oCmd := exec.Command("sh", []string{"-c", cmd}...)
  oCmd.Dir = dir
	return oCmd, nil
}
